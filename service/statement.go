package service

import (
	"context"

	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils"
	"github.com/takeme-id/core/utils/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func WithdrawFeeStatement(balanceID primitive.ObjectID, time string, transactionCode string,
	amount int) domain.Statement {
	return domain.Statement{
		BalanceID:   balanceID,
		Time:        time,
		Description: "Withdraw for fee from " + transactionCode,
		Reference:   transactionCode,
		Withdraw:    amount,
		Deposit:     0,
		Type:        domain.STATEMENT_TYPE_FEE,
	}
}

func DepositFeeStatement(balanceID primitive.ObjectID, time string, transactionCode string,
	amount int) domain.Statement {
	return domain.Statement{
		BalanceID:   balanceID,
		Time:        time,
		Description: "Deposit for fee from " + transactionCode,
		Reference:   transactionCode,
		Withdraw:    0,
		Deposit:     amount,
		Type:        domain.STATEMENT_TYPE_FEE,
	}
}

func WithdrawTransactionStatement(balanceID primitive.ObjectID, time string, transactionCode string,
	amount int) domain.Statement {
	return domain.Statement{
		BalanceID:   balanceID,
		Time:        time,
		Description: "Withdraw for " + transactionCode,
		Reference:   transactionCode,
		Withdraw:    amount,
		Deposit:     0,
		Type:        domain.STATEMENT_TYPE_TRANSACTION,
	}
}

func DepositTransactionStatement(balanceID primitive.ObjectID, time string, transactionCode string,
	amount int) domain.Statement {
	return domain.Statement{
		BalanceID:   balanceID,
		Time:        time,
		Description: "Deposit for " + transactionCode,
		Reference:   transactionCode,
		Withdraw:    0,
		Deposit:     amount,
		Type:        domain.STATEMENT_TYPE_TRANSACTION,
	}
}

func StatementsByBalanceID(balanceID primitive.ObjectID, page string, limit string) ([]domain.Statement, error) {
	query := bson.M{"balance_id": balanceID}

	var results []domain.Statement
	cursor, err := database.Find(domain.STATEMENT_COLLECTION_NAME, query, page, limit)
	err = cursor.All(context.TODO(), &results)

	if err != nil {
		return []domain.Statement{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	return results, nil
}

func StatementSaveOne(model domain.Statement, session mongo.SessionContext) error {
	err := database.SessionSaveOne(model, session)
	if err != nil {
		return err
	}

	return nil
}
