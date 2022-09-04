package acceptpayment

import (
	"os"
	"strconv"
	"time"

	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/service"
	"github.com/takeme-id/core/usecase"
	"github.com/takeme-id/core/usecase/transaction"
	"github.com/takeme-id/core/utils"
	"github.com/takeme-id/core/utils/gateway"
)

type AcceptCard struct {
	corporate          domain.Corporate
	from               domain.Card
	to                 domain.TransactionObject
	balance            domain.Balance
	amount             int
	currency           string
	reference          string
	transactionUsecase transaction.Base
}

func (self AcceptCard) Initialize(from domain.Card, balanceID string, amount int,
	reference string, currency string, returnURL string, externalID string) (string, string, error) {

	gateway := gateway.StripeGateway{}

	status, authURL, err := gateway.ChargeCard(balanceID, amount, returnURL, from, externalID)
	if err != nil {
		return "", "", err
	}

	return status, authURL, nil
}

func (self AcceptCard) InitializeSubscribe(from domain.Card, balanceID string, amount int,
	reference string, currency string, returnURL string, externalID string, interval string) (string, string, string, error) {

	gateway := gateway.StripeGateway{}

	status, authURL, subsID, err := gateway.ChargeCardSubscribe(balanceID, amount, returnURL, from, externalID, interval)
	if err != nil {
		return "", "", subsID, err
	}

	return status, authURL, subsID, nil
}

func (self AcceptCard) Execute(from domain.Card, balanceID string, amount int,
	reference string, currency string, externalID string) (domain.Transaction, domain.Balance, error) {

	balance, owner, corporate, err := identifyBalance(balanceID)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	gateway := gateway.StripeGateway{}
	self.corporate = corporate
	self.from = from
	self.to = owner
	self.balance = balance
	self.amount = amount
	self.reference = reference
	self.transactionUsecase = transaction.Base{}
	self.currency = currency

	var statements []domain.Statement

	transaction, transactionStatement, err := createTransaction(self.corporate, self.balance, self.from,
		self.to, self.amount, self.reference, gateway, externalID)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(corporate, balance, transaction)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	statements = append(statements, transactionStatement)
	statements = append(statements, feeStatement...)

	err = validateCurrency(self.currency, balance)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	err = self.transactionUsecase.Commit(statements, &transaction)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	go usecase.PublishAcceptPaymentCallback(corporate, balance, transaction)

	return transaction, balance, nil
}

func identifyBalance(balanceID string) (domain.Balance, domain.TransactionObject, domain.Corporate, error) {
	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil {
		return domain.Balance{}, domain.TransactionObject{}, domain.Corporate{},
			utils.ErrorBadRequest(utils.InvalidBalanceID, "Balance id not found")
	}

	var balanceOwner domain.TransactionObject
	ownerID := balance.Owner.ID.Hex()
	corporateID := balance.CorporateID.Hex()
	balanceOwnerType := balance.Owner.Type

	if balanceOwnerType == domain.ACTOR_TYPE_CORPORATE {
		corporate, err := service.CorporateByIDNoSession(ownerID)
		if err != nil {
			return domain.Balance{}, domain.TransactionObject{}, domain.Corporate{},
				utils.ErrorBadRequest(utils.CorporateNotFound, "Corporate id not found")
		}

		balanceOwner = corporate.ToTransactionObject()
	} else {
		user, err := service.UserByIDNoSession(ownerID)
		if err != nil {
			return domain.Balance{}, domain.TransactionObject{}, domain.Corporate{},
				utils.ErrorBadRequest(utils.UserNotFound, "User id not found")
		}

		balanceOwner = user.ToTransactionObject()
	}

	corporate, err := service.CorporateByIDNoSession(corporateID)
	if err != nil {
		return domain.Balance{}, domain.TransactionObject{}, domain.Corporate{},
			utils.ErrorBadRequest(utils.CorporateNotFound, "Corporate id not found")
	}

	return balance, balanceOwner, corporate, nil
}

func createTransaction(corporate domain.Corporate, balance domain.Balance, from domain.Card,
	to domain.TransactionObject, subAmount int, reference string, gateway gateway.Gateway, externalID string) (domain.Transaction, domain.Statement, error) {

	var totalFee = 0
	if balance.Owner.Type == domain.ACTOR_TYPE_USER {

		s, err := strconv.ParseFloat(corporate.FeeUser.AcceptPaymentCard, 64)
		if err != nil {
			return domain.Transaction{}, domain.Statement{}, utils.ErrorBadRequest(utils.WrongAcceptCardFee, "Cannot convert accept payment card fee")
		}

		totalFee = int(float64(subAmount) * s)
	} else {
		s, err := strconv.ParseFloat(corporate.FeeCorporate.AcceptPaymentCard, 64)
		if err != nil {
			return domain.Transaction{}, domain.Statement{}, utils.ErrorBadRequest(utils.WrongAcceptCardFee, "Cannot convert accept payment card fee")
		}

		totalFee = int(float64(subAmount) * s)
	}

	transcation := domain.Transaction{
		TransactionCode:  utils.GenerateTransactionCode("2"),
		CorporateID:      corporate.ID,
		Type:             domain.ACCEPT_PAYMENT_CARD,
		Method:           domain.METHOD_CARD,
		ToBalanceID:      balance.ID,
		FromBalanceID:    balance.ID,
		From:             from.ToTransactionObject(),
		To:               to,
		TotalFee:         totalFee,
		SubAmount:        subAmount,
		Amount:           subAmount - totalFee,
		Time:             time.Now().Format(os.Getenv("TIME_FORMAT")),
		Notes:            "",
		Status:           domain.COMPLETED_STATUS,
		Unpaid:           false,
		ExternalID:       externalID,
		Gateway:          gateway.Name(),
		GatewayReference: reference,
		Currency:         corporate.Currency,
	}

	statement := service.DepositTransactionStatement(
		balance.ID, transcation.Time, transcation.TransactionCode, subAmount)

	return transcation, statement, nil
}
