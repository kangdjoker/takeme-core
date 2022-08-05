package service

import (
	"context"
	"net/http"

	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/domain/dto"
	"github.com/takeme-id/core/utils"
	"github.com/takeme-id/core/utils/database"
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
	corporateID := r.Header.Get("corporate")

	model := domain.Corporate{}
	cursor := database.FindOneByID(domain.CORPORATE_COLLECTION, corporateID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Corporate{}, utils.ErrorBadRequest(utils.InvalidCorporateKey, "Corporate not found")
	}

	err = ValidateCorporateExist(model)
	if err != nil {
		return domain.Corporate{}, err
	}

	err = ValidateCorporateLocked(model)
	if err != nil {
		return domain.Corporate{}, err
	}

	return model, nil
}

func CorporateUpdateOne(model *domain.Corporate, session mongo.SessionContext) error {
	err := database.SessionUpdateOne(model, session)
	if err != nil {
		return err
	}

	return nil
}

func CorporateReduceAccessAttempt(corporateID string, session mongo.SessionContext) error {

	corporate, err := CorporateByID(corporateID, session)
	if err != nil {
		return err
	}

	err = ValidateCorporateLocked(corporate)
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

	err = CorporateUpdateOne(&corporate, session)
	if err != nil {
		return err
	}

	return nil
}

func CorporateSavePIN(corporate *domain.Corporate, pin string, session mongo.SessionContext) error {

	pin, err := utils.RSADecrypt(pin)
	if err != nil {
		return utils.ErrorInternalServer(utils.DecryptError, err.Error())
	}

	corporate.PIN = pin

	err = CorporateUpdateOne(corporate, session)
	if err != nil {
		return err
	}

	return nil
}

func CorporateChangeNewPIN(corporate *domain.Corporate, newPIN string, session mongo.SessionContext) error {
	corporate.PIN = newPIN

	err := CorporateUpdateOne(corporate, session)
	if err != nil {
		return err
	}

	return nil
}

func ValidateCorporateLocked(corporate domain.Corporate) error {
	if corporate.Active == false {
		return utils.ErrorBadRequest(utils.CorporateLocked, "Corporate Locked")
	}

	return nil
}

func ValidateCorporateExist(corporate domain.Corporate) error {
	if corporate.Name == "" {
		return utils.ErrorBadRequest(utils.InvalidCorporateKey, "Corporate not found")
	}

	return nil
}

func CorporateDTOByID(corporateID string) (dto.Corporate, error) {

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
	cursor, err := database.Aggregate(domain.CORPORATE_COLLECTION, query)
	cursor.All(context.TODO(), &result)
	if err != nil {
		return dto.Corporate{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed or cannot decode")
	}

	return result[0], nil
}
