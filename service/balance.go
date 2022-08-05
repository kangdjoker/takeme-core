package service

import (
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func BalanceInitialization(id primitive.ObjectID, corporateID primitive.ObjectID, owner domain.ActorObject,
	name string, session mongo.SessionContext) (domain.Balance, error) {

	model := domain.Balance{
		ID:          id,
		CorporateID: corporateID,
		Owner:       owner,
		Name:        name,
		Amount:      0,
	}

	err := BalanceSaveOne(&model, session)
	if err != nil {
		return model, err
	}

	return model, nil
}

func BalanceCreate(corporateID primitive.ObjectID, owner domain.ActorObject,
	name string, session mongo.SessionContext) (domain.Balance, error) {

	model := domain.Balance{
		CorporateID: corporateID,
		Owner:       owner,
		Name:        name,
		Amount:      0,
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

func BalanceSaveOneNoSession(model *domain.Balance) error {
	err := database.SaveOne(domain.BALANCE_COLLECTION, model)
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

func BalanceUpdate(model domain.Balance, session mongo.SessionContext) error {
	err := database.SessionUpdateOne(&model, session)
	if err != nil {
		return err
	}

	return nil
}
