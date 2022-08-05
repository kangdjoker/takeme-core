package service

import (
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func FraudSave(fraud domain.Fraud, session mongo.SessionContext) error {
	err := database.SessionSaveOne(&fraud, session)
	if err != nil {
		return err
	}

	return nil
}
