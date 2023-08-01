package usecase

import (
	"context"
	"fmt"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func LogInformation(data interface{}) (domain.Log, error) {
	var log domain.Log

	createLog := func(session mongo.SessionContext) error {

		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}
		log, err = service.LogCreate(data, session)

		if err != nil {
			session.AbortTransaction(session)
			return err
		}
		return database.CommitWithRetry(session)
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, createLog)
		},
	)

	if err != nil {
		return log, utils.ErrorInternalServer(utils.DBStartTransactionFailed,
			fmt.Sprintf("Failed initialize log (%v)", log.GetDocumentID()))
	}

	return log, nil
}
