package security

import (
	"context"
	"fmt"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func InvalidCorporateAuth(paramLog *basic.ParamLog, corporate domain.Corporate) {

	transactionFunction := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Invalid corporate auth start transaction")
		}

		err = service.CorporateReduceAccessAttempt(paramLog, corporate.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)

	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, transactionFunction)
		},
	)

	if err != nil {
		basic.LogError(paramLog, fmt.Sprintf("Invalid corporate auth failed because %v ", err.Error()))
	}
}

func InvalidUserAuth(paramLog *basic.ParamLog, user domain.User) {
	transactionFunction := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Invalid corporate auth start transaction")
		}

		err = service.UserReduceAccessAttempt(paramLog, user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		fraud := domain.CreateFraud(domain.USER_FAILED_ATTEMPT, user, domain.USER_COLLECTION)
		err = service.FraudSave(fraud, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)

	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, transactionFunction)
		},
	)

	if err != nil {
		basic.LogError(paramLog, fmt.Sprintf("Invalid corporate auth failed because %v ", err.Error()))
	}
}

func LockUser(paramLog *basic.ParamLog, userID string) {

	transactionFunction := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Invalid lock user start transaction")
		}

		user, err := service.UserByID(paramLog, userID, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.UserLock(paramLog, &user, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)

	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, transactionFunction)
		},
	)

	if err != nil {
		basic.LogError(paramLog, fmt.Sprintf("Failed lock user with id  %v ", userID))
	}
}
