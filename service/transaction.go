package service

import (
	"context"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TransactionSaveOne(model *domain.Transaction, session mongo.SessionContext) error {
	err := database.SessionSaveOne(model, session)
	if err != nil {
		return err
	}

	return nil
}

func TransactionPendingByCodeNoSession(code string) (domain.Transaction, error) {
	var transaction domain.Transaction
	query := bson.M{"transaction_code": code, "status": "Pending"}
	cursor := database.FindOne(domain.TRANSACTION_COLLECTION, query)
	err := cursor.Decode(&transaction)
	if err != nil {
		return domain.Transaction{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	if transaction.TransactionCode == "" {
		return domain.Transaction{}, utils.ErrorBadRequest(utils.TransactionNotFound, "Transaction not found")
	}

	return transaction, nil
}

func TransactionPendingByReferenceNoSession(code string) (domain.Transaction, error) {
	var transaction domain.Transaction
	query := bson.M{"gateway_reference": code, "status": "Pending"}
	cursor := database.FindOne(domain.TRANSACTION_COLLECTION, query)
	err := cursor.Decode(&transaction)
	if err != nil {
		return domain.Transaction{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	if transaction.TransactionCode == "" {
		return domain.Transaction{}, utils.ErrorBadRequest(utils.TransactionNotFound, "Transaction not found")
	}

	return transaction, nil
}

func TransactionByGatewayReferenceNoSession(code string) (domain.Transaction, error) {
	var transaction domain.Transaction
	query := bson.M{"gateway_reference": code}
	cursor := database.FindOne(domain.TRANSACTION_COLLECTION, query)
	err := cursor.Decode(&transaction)
	if err != nil {
		return domain.Transaction{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	if transaction.TransactionCode == "" {
		return domain.Transaction{}, utils.ErrorBadRequest(utils.TransactionNotFound, "Transaction not found")
	}

	return transaction, nil
}

func TransactionByCodeNoSession(code string) (domain.Transaction, error) {
	var transaction domain.Transaction
	query := bson.M{"transaction_code": code}
	cursor := database.FindOne(domain.TRANSACTION_COLLECTION, query)
	err := cursor.Decode(&transaction)
	if err != nil {
		return domain.Transaction{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	if transaction.TransactionCode == "" {
		return domain.Transaction{}, utils.ErrorBadRequest(utils.TransactionNotFound, "Transaction not found")
	}

	return transaction, nil
}

func TransactionByID(ID string, session mongo.SessionContext) (domain.Transaction, error) {
	model := domain.Transaction{}
	cursor := database.SessionFindOneByID(domain.TRANSACTION_COLLECTION, ID, session)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.Transaction{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed or cannot decode")
	}

	return model, nil
}

func TransactionPendingByCode(code string, session mongo.SessionContext) (domain.Transaction, error) {
	var transaction domain.Transaction
	query := bson.M{"transaction_code": code, "status": "Pending"}
	cursor := database.SessionFindOne(domain.TRANSACTION_COLLECTION, query, session)
	err := cursor.Decode(&transaction)
	if err != nil {
		return domain.Transaction{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	if transaction.TransactionCode == "" {
		return domain.Transaction{}, utils.ErrorBadRequest(utils.TransactionNotFound, "Transaction not found")
	}

	return transaction, nil
}

func TransactionUpdateOne(model *domain.Transaction, session mongo.SessionContext) error {
	err := database.SessionUpdateOne(model, session)
	if err != nil {
		return err
	}

	return nil
}

func TransactionsByActorNoSession(accountNumber string, ownBalance []primitive.ObjectID, page string, limit string) ([]domain.Transaction, error) {

	query := bson.M{}
	orQuery := []bson.M{}
	orQuery = append(orQuery, bson.M{"from.account_number": accountNumber}, bson.M{"to.account_number": accountNumber})

	for _, a := range ownBalance {
		orQuery = append(orQuery, bson.M{"balance_id": a})
	}

	query["$or"] = orQuery

	var transaction []domain.Transaction
	cursor, err := database.Find(domain.TRANSACTION_COLLECTION, query, page, limit)
	err = cursor.All(context.TODO(), &transaction)

	if err != nil {
		return []domain.Transaction{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed")
	}

	return transaction, nil
}
