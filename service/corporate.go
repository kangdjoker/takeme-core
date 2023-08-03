package service

import (
	"context"
	"net/http"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/domain/dto"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CorporateSave(corporate domain.Corporate, session mongo.SessionContext) error {
	err := database.SessionSaveOne(&corporate, session)
	if err != nil {
		return err
	}

	return nil
}

func CorporateByIDNoSession(ID string) (domain.Corporate, error) {
	model := domain.Corporate{}
	cursor := database.FindOneByID(domain.CORPORATE_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Corporate{}, err
	}

	return model, nil
}

func CorporateByID(ID string, session mongo.SessionContext) (domain.Corporate, error) {
	model := domain.Corporate{}
	cursor := database.SessionFindOneByID(domain.CORPORATE_COLLECTION, ID, session)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Corporate{}, err
	}

	return model, nil
}

func CorporateByRequest(r *http.Request) (domain.Corporate, error) {
	ioCloser, span, tag := basic.RequestToTracing(r)
	paramLog := &basic.ParamLog{Span: span, TrCloser: ioCloser, Tag: tag}
	corporateID := r.Header.Get("corporate")

	model := domain.Corporate{}
	cursor := database.FindOneByID(domain.CORPORATE_COLLECTION, corporateID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Corporate{}, utils.ErrorBadRequest(paramLog, utils.InvalidCorporateKey, "Corporate not found")
	}

	err = ValidateCorporateExist(paramLog, model)
	if err != nil {
		return domain.Corporate{}, err
	}

	err = ValidateCorporateLocked(paramLog, model)
	if err != nil {
		return domain.Corporate{}, err
	}

	return model, nil
}

func CorporateUpdateOne(paramLog *basic.ParamLog, model *domain.Corporate, session mongo.SessionContext) error {
	err := database.SessionUpdateOne(paramLog, model, session)
	if err != nil {
		return err
	}

	return nil
}

func CorporateReduceAccessAttempt(paramLog *basic.ParamLog, corporateID string, session mongo.SessionContext) error {

	corporate, err := CorporateByID(corporateID, session)
	if err != nil {
		return err
	}

	err = ValidateCorporateLocked(paramLog, corporate)
	if err != nil {
		return err
	}

	remaining := corporate.AccessAttempt
	if remaining <= 0 {
		corporate.AccessAttempt = 0
		corporate.Active = false
	} else {
		corporate.AccessAttempt -= 1
	}

	err = CorporateUpdateOne(paramLog, &corporate, session)
	if err != nil {
		return err
	}

	return nil
}

func CorporateSavePIN(paramLog *basic.ParamLog, corporate *domain.Corporate, pin string, session mongo.SessionContext) error {

	pin, err := utils.RSADecrypt(pin)
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.DecryptError, err.Error())
	}

	corporate.PIN = pin

	err = CorporateUpdateOne(paramLog, corporate, session)
	if err != nil {
		return err
	}

	return nil
}

func CorporateChangeNewPIN(paramLog *basic.ParamLog, corporate *domain.Corporate, newPIN string, session mongo.SessionContext) error {
	corporate.PIN = newPIN

	err := CorporateUpdateOne(paramLog, corporate, session)
	if err != nil {
		return err
	}

	return nil
}

func ValidateCorporateLocked(paramLog *basic.ParamLog, corporate domain.Corporate) error {
	if corporate.Active == false {
		return utils.ErrorBadRequest(paramLog, utils.CorporateLocked, "Corporate Locked")
	}

	return nil
}

func ValidateCorporateExist(paramLog *basic.ParamLog, corporate domain.Corporate) error {
	if corporate.Name == "" {
		return utils.ErrorBadRequest(paramLog, utils.InvalidCorporateKey, "Corporate not found")
	}

	return nil
}

func CorporateDTOByID(paramLog *basic.ParamLog, corporateID string) (dto.Corporate, error) {

	objectID, err := primitive.ObjectIDFromHex(corporateID)

	query := []bson.M{
		{"$match": bson.M{"_id": objectID}},
		{"$unwind": "$list_balance"},
		{
			"$lookup": bson.M{
				"from":         "balance",
				"localField":   "list_balance.balance_id",
				"foreignField": "_id",
				"as":           "list_balance.detail",
			},
		},
		{"$unwind": "$list_balance.detail"},
		{
			"$group": bson.M{
				"_id": "$_id",
				"list_balance": bson.M{
					"$push": "$list_balance",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "corporate",
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "corporate_detail",
			},
		},
		{"$unwind": "$corporate_detail"},
		{
			"$addFields": bson.M{
				"corporate_detail.list_balance": "$list_balance",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$corporate_detail",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "balance",
				"localField":   "main_balance",
				"foreignField": "_id",
				"as":           "main_balance",
			},
		},
		{"$unwind": "$main_balance"},
	}

	var result []dto.Corporate
	cursor, err := database.Aggregate(paramLog, domain.CORPORATE_COLLECTION, query)
	cursor.All(context.TODO(), &result)
	if err != nil {
		return dto.Corporate{}, utils.ErrorInternalServer(paramLog, utils.QueryFailed, "Query failed or cannot decode")
	}

	return result[0], nil
}
