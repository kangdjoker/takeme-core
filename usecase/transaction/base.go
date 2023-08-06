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
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Base struct {
}

func (self Base) CreateFeeStatement(paramLog *basic.ParamLog, corporate domain.Corporate, balance domain.Balance,
	transaction domain.Transaction) ([]domain.Statement, error) {
	feeCalculator := usecase.CalculateFee{}
	feeCalculator.Initialize(corporate, balance, transaction)

	statements, err := feeCalculator.CalculateByOwnerAndTransaction(paramLog)
	if err != nil {
		return []domain.Statement{}, err
	}

	return statements, nil
}

func (self Base) RollbackFeeStatement(paramLog *basic.ParamLog, corporate domain.Corporate, balance domain.Balance,
	transaction domain.Transaction) ([]domain.Statement, error) {
	feeCalculator := usecase.CalculateFee{}
	feeCalculator.Initialize(corporate, balance, transaction)

	feeStatements, err := feeCalculator.CalculateByOwnerAndTransaction(paramLog)
	statements := feeCalculator.RollbackFeeStatement(feeStatements)

	if err != nil {
		return []domain.Statement{}, err
	}

	return statements, nil
}

func (self Base) Commit(paramLog *basic.ParamLog, statements []domain.Statement, transaction *domain.Transaction) error {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			basic.LogInformation(paramLog, "Error: "+err.Error())
			session.AbortTransaction(session)
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		err = service.TransactionSaveOne(transaction, session)
		if err != nil {
			basic.LogInformation(paramLog, "Error: "+err.Error())
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

		err = adjustBalanceWithStatement(paramLog, statements, session)
		if err != nil {
			basic.LogInformation(paramLog, "Error: "+err.Error())
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

func (self Base) UpdatingTransactionDetail(paramLog *basic.ParamLog, transaction *domain.Transaction) error {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			basic.LogError(paramLog, "Error: "+err.Error())
			session.AbortTransaction(session)
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		err = service.TransactionUpdateOne(paramLog, transaction, session)
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

func (self Base) CommitRollback(paramLog *basic.ParamLog, statements []domain.Statement) error {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		basic.LogInformation(paramLog, "StartTransaction")

		if err != nil {
			session.AbortTransaction(session)
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}
		basic.LogInformation(paramLog, "adjustBalanceWithStatement")

		err = adjustBalanceWithStatement(paramLog, statements, session)
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

func adjustBalanceWithStatement(paramLog *basic.ParamLog, statements []domain.Statement, session mongo.SessionContext) error {

	for _, statement := range statements {
		if statement.Deposit != 0 {
			basic.LogInformation2(paramLog, "adjustBalanceWithStatement.Deposit", statement)
			err := usecase.DepositBalance(paramLog, statement, session)
			if err != nil {
				return err
			}
		} else if statement.Withdraw != 0 {
			basic.LogInformation2(paramLog, "adjustBalanceWithStatement.Withdraw", statement)
			err := usecase.WithdrawBalance(paramLog, statement, session)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
