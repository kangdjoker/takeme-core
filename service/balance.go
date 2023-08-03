package service

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func BalanceInitialization(id primitive.ObjectID, corporateID primitive.ObjectID, owner domain.ActorObject,
	name string, currency string, session mongo.SessionContext) (domain.Balance, error) {

	model := domain.Balance{
		ID:          id,
		CorporateID: corporateID,
		Owner:       owner,
		Name:        name,
		Amount:      0,
		Currency:    currency,
	}

	err := BalanceSaveOne(&model, session)
	if err != nil {
		return model, err
	}

	return model, nil
}

func BalanceCreate(corporateID primitive.ObjectID, owner domain.ActorObject,
	name string, currency string, session mongo.SessionContext) (domain.Balance, error) {

	model := domain.Balance{
		CorporateID: corporateID,
		Owner:       owner,
		Name:        name,
		Amount:      0,
		Currency:    currency,
	}

	err := BalanceSaveOne(&model, session)
	if err != nil {
		return model, err
	}

	return model, nil
}

func BalanceSaveOne(model *domain.Balance, session mongo.SessionContext) error {
	err := database.SessionSaveOne(model, session)
	if err != nil {
		return err
	}

	return nil
}

func BalanceSaveOneNoSession(paramLog *basic.ParamLog, model *domain.Balance) error {
	err := database.SaveOne(paramLog, domain.BALANCE_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}

func BalanceByID(ID string, session mongo.SessionContext) (domain.Balance, error) {
	model := domain.Balance{}
	cursor := database.SessionFindOneByID(domain.BALANCE_COLLECTION, ID, session)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Balance{}, err
	}

	return model, nil
}

func BalanceByIDNoSession(ID string) (domain.Balance, error) {
	model := domain.Balance{}
	cursor := database.FindOneByID(domain.BALANCE_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Balance{}, err
	}

	return model, nil
}

func BalanceUpdate(paramLog *basic.ParamLog, model domain.Balance, session mongo.SessionContext) error {
	err := database.SessionUpdateOne(paramLog, &model, session)
	if err != nil {
		return err
	}

	return nil
}
