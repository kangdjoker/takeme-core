package service

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func FraudSave(fraud domain.Fraud, session mongo.SessionContext) error {
	err := database.SessionSaveOne(&fraud, session)
	if err != nil {
		return err
	}

	return nil
}
