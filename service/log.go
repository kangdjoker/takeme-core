package service

import (
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func LogInitialization(id primitive.ObjectID, data interface{}, session mongo.SessionContext) (domain.Log, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	model := domain.Log{
		ID:         id,
		Data:       data,
		TimeCreate: now,
	}

	err := LogSaveOne(&model, session)
	if err != nil {
		return model, err
	}

	return model, nil
}

func LogCreate(data interface{}, session mongo.SessionContext) (domain.Log, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	model := domain.Log{
		Data:       data,
		TimeCreate: now,
	}

	err := LogSaveOne(&model, session)
	if err != nil {
		return model, err
	}

	return model, nil
}

func LogSaveOne(model *domain.Log, session mongo.SessionContext) error {
	err := database.SessionSaveOne(model, session)
	if err != nil {
		return err
	}

	return nil
}

func LogSaveOneNoSession(model *domain.Log) error {
	err := database.SaveOne(domain.LOG_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}

func LogById(ID string, session mongo.SessionContext) (domain.Log, error) {
	model := domain.Log{}
	cursor := database.SessionFindOneByID(domain.LOG_COLLECTION, ID, session)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Log{}, err
	}

	return model, nil
}

func LogByIDNoSession(ID string) (domain.Log, error) {
	model := domain.Log{}
	cursor := database.FindOneByID(domain.LOG_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Log{}, err
	}

	return model, nil
}

func LogUpdate(model domain.Log, session mongo.SessionContext) error {
	err := database.SessionUpdateOne(&model, session)
	if err != nil {
		return err
	}

	return nil
}
