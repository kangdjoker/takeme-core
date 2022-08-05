package service

import (
	"context"
	"os"
	"time"

	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils"
	"github.com/takeme-id/core/utils/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateRAB(corporate domain.Corporate, balance domain.Balance, requester domain.ActorObject, owner domain.ActorObject,
	access string) (domain.RequestAccessBalance, error) {
	model := domain.RequestAccessBalance{
		CorporateID:      corporate.ID,
		BalanceID:        balance.ID,
		Time:             time.Now().Format(os.Getenv("TIME_FORMAT")),
		BalanceRequester: requester,
		BalanceOwner:     owner,
		Access:           access,
		Status:           domain.REQUEST_ACCESS_BALANCE_STATUS_PENDING,
	}

	err := RABSaveOne(&model)
	if err != nil {
		return domain.RequestAccessBalance{}, err
	}

	return model, nil
}

func RABSaveOne(model *domain.RequestAccessBalance) error {
	err := database.SaveOne(domain.RAB_COLLECTION_NAME, model)
	if err != nil {
		return err
	}

	return nil
}

func RABByID(ID string) (domain.RequestAccessBalance, error) {
	model := domain.RequestAccessBalance{}
	cursor := database.FindOneByID(domain.RAB_COLLECTION_NAME, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.RequestAccessBalance{},
			utils.ErrorInternalServer(utils.QueryFailed, "Query failed or cannot decode")
	}

	return model, nil
}

func RABUpdateOne(model *domain.RequestAccessBalance) error {
	err := database.UpdateOne(domain.RAB_COLLECTION_NAME, model)
	if err != nil {
		return err
	}

	return nil
}

func RABByRequsterID(ID string, status string) ([]domain.RequestAccessBalance, error) {
	objectID, _ := primitive.ObjectIDFromHex(ID)
	query := bson.M{"balance_requester._id": objectID, "status": bson.M{"$regex": status, "$options": "i"}}

	var models []domain.RequestAccessBalance
	cursor, err := database.Find(domain.RAB_COLLECTION_NAME, query, "1", "1000")
	err = cursor.All(context.TODO(), &models)

	if err != nil {
		return []domain.RequestAccessBalance{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	return models, nil
}

func RABByOwnerID(ID string, status string) ([]domain.RequestAccessBalance, error) {
	objectID, _ := primitive.ObjectIDFromHex(ID)
	query := bson.M{"balance_owner._id": objectID, "status": bson.M{"$regex": status, "$options": "i"}}

	var models []domain.RequestAccessBalance
	cursor, err := database.Find(domain.RAB_COLLECTION_NAME, query, "1", "1000")
	err = cursor.All(context.TODO(), &models)

	if err != nil {
		return []domain.RequestAccessBalance{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	return models, nil
}
