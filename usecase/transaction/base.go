package transaction

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/database"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Base struct {
}

func (self Base) CreateFeeStatement(corporate domain.Corporate, balance domain.Balance,
	transaction domain.Transaction) ([]domain.Statement, error) {
	feeCalculator := usecase.CalculateFee{}
	feeCalculator.Initialize(corporate, balance, transaction)

	statements, err := feeCalculator.CalculateByOwnerAndTransaction()
	if err != nil {
		return []domain.Statement{}, err
	}

	return statements, nil
}

func (self Base) RollbackFeeStatement(corporate domain.Corporate, balance domain.Balance,
	transaction domain.Transaction) ([]domain.Statement, error) {
	feeCalculator := usecase.CalculateFee{}
	feeCalculator.Initialize(corporate, balance, transaction)

	feeStatements, err := feeCalculator.CalculateByOwnerAndTransaction()
	statements := feeCalculator.RollbackFeeStatement(feeStatements)

	if err != nil {
		return []domain.Statement{}, err
	}

	return statements, nil
}

func (self Base) Commit(tag string, statements []domain.Statement, transaction *domain.Transaction) error {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			logrus.Info("Error: " + err.Error())
			session.AbortTransaction(session)
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		err = service.TransactionSaveOne(transaction, session)
		if err != nil {
			logrus.Info("Error: " + err.Error())
			session.AbortTransaction(session)
			if strings.Contains(err.Error(), "E11000") {
				return utils.CustomError{
					HttpStatus:  http.StatusBadRequest,
					Code:        11000,
					Description: "Duplicate Request ID",
					Time:        time.Now().Format(os.Getenv("TIME_FORMAT")),
				}
			} else {
				return err
			}
		}

		err = adjustBalanceWithStatement(tag, statements, session)
		if err != nil {
			logrus.Info("Error: " + err.Error())
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)

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

func (self Base) UpdatingTransactionDetail(transaction *domain.Transaction) error {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			logrus.Info("Error: " + err.Error())
			session.AbortTransaction(session)
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		err = service.TransactionUpdateOne(transaction, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)

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

func (self Base) CommitRollback(tag string, statements []domain.Statement) error {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			session.AbortTransaction(session)
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		err = adjustBalanceWithStatement(tag, statements, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)

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

func adjustBalanceWithStatement(tag string, statements []domain.Statement, session mongo.SessionContext) error {

	for _, statement := range statements {
		if statement.Deposit != 0 {
			err := usecase.DepositBalance(statement, session)
			if err != nil {
				return err
			}
		} else if statement.Withdraw != 0 {
			err := usecase.WithdrawBalance(tag, statement, session)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
