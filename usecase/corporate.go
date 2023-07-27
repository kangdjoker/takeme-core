package usecase

import (
	"context"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/domain/dto"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func CorporateSavePIN(corporate domain.Corporate, encryptedPIN string) error {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		corporate, err = service.CorporateByID(corporate.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.CorporateSavePIN(&corporate, encryptedPIN, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, function)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func CorporateChangePIN(corporate domain.Corporate, encryptedOldPIN string, encryptedNewPIN string) error {

	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		newPIN, err := utils.RSADecrypt(encryptedNewPIN)
		if err != nil {
			return err
		}

		err = ValidateFormatPIN(newPIN)
		if err != nil {
			return err
		}

		corporate, err = service.CorporateByID(corporate.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = ValidateActorPIN(corporate, encryptedOldPIN)
		if err != nil {
			return err
		}

		err = service.CorporateChangeNewPIN(&corporate, newPIN, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, function)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func CorporateCheck(corporate domain.Corporate) (dto.Corporate, error) {

	result, err := service.CorporateDTOByID(corporate.ID.Hex())
	if err != nil {
		return dto.Corporate{}, err
	}

	return result, nil
}
