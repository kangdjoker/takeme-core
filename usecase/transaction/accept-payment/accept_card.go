package acceptpayment

import (
	"os"
	"strconv"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/gateway"
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

func (self AcceptCard) Initialize(paramLog *basic.ParamLog, from domain.Card, balanceID string, amount int,
	reference string, currency string, returnURL string, externalID string) (string, string, error) {

	gateway := gateway.StripeGateway{}

	status, authURL, err := gateway.ChargeCard(paramLog, balanceID, amount, returnURL, from, externalID)
	if err != nil {
		return "", "", err
	}

	return status, authURL, nil
}

func (self AcceptCard) InitializeSubscribe(paramLog *basic.ParamLog, from domain.Card, balanceID string, amount int,
	reference string, currency string, returnURL string, externalID string, interval string) (string, string, string, error) {

	gateway := gateway.StripeGateway{}

	status, authURL, subsID, err := gateway.ChargeCardSubscribe(paramLog, balanceID, amount, returnURL, from, externalID, interval)
	if err != nil {
		return "", "", subsID, err
	}

	return status, authURL, subsID, nil
}

func (self AcceptCard) Execute(paramLog *basic.ParamLog, from domain.Card, balanceID string, amount int,
	reference string, currency string, externalID string, requestId string) (domain.Transaction, domain.Balance, error) {

	balance, owner, corporate, err := identifyBalance(paramLog, balanceID)
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

	transaction, transactionStatement, err := createTransaction(paramLog, self.corporate, self.balance, self.from,
		self.to, self.amount, self.reference, gateway, externalID, requestId)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(paramLog, corporate, balance, transaction)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	statements = append(statements, transactionStatement)
	statements = append(statements, feeStatement...)

	err = validateCurrency(paramLog, self.currency, balance)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	err = self.transactionUsecase.Commit(paramLog, statements, &transaction)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	go usecase.PublishAcceptPaymentCallback(paramLog, corporate, balance, transaction)

	return transaction, balance, nil
}

func identifyBalance(paramLog *basic.ParamLog, balanceID string) (domain.Balance, domain.TransactionObject, domain.Corporate, error) {
	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil {
		return domain.Balance{}, domain.TransactionObject{}, domain.Corporate{},
			utils.ErrorBadRequest(paramLog, utils.InvalidBalanceID, "Balance id not found")
	}

	var balanceOwner domain.TransactionObject
	ownerID := balance.Owner.ID.Hex()
	corporateID := balance.CorporateID.Hex()
	balanceOwnerType := balance.Owner.Type

	if balanceOwnerType == domain.ACTOR_TYPE_CORPORATE {
		corporate, err := service.CorporateByIDNoSession(ownerID)
		if err != nil {
			return domain.Balance{}, domain.TransactionObject{}, domain.Corporate{},
				utils.ErrorBadRequest(paramLog, utils.CorporateNotFound, "Corporate id not found")
		}

		balanceOwner = corporate.ToTransactionObject()
	} else {
		user, err := service.UserByIDNoSession(paramLog, ownerID)
		if err != nil {
			return domain.Balance{}, domain.TransactionObject{}, domain.Corporate{},
				utils.ErrorBadRequest(paramLog, utils.UserNotFound, "User id not found")
		}

		balanceOwner = user.ToTransactionObject()
	}

	corporate, err := service.CorporateByIDNoSession(corporateID)
	if err != nil {
		return domain.Balance{}, domain.TransactionObject{}, domain.Corporate{},
			utils.ErrorBadRequest(paramLog, utils.CorporateNotFound, "Corporate id not found")
	}

	return balance, balanceOwner, corporate, nil
}

func createTransaction(paramLog *basic.ParamLog, corporate domain.Corporate, balance domain.Balance, from domain.Card,
	to domain.TransactionObject, subAmount int, reference string, gateway gateway.Gateway, externalID string, requestId string) (domain.Transaction, domain.Statement, error) {

	var totalFee = 0
	if balance.Owner.Type == domain.ACTOR_TYPE_USER {

		s, err := strconv.ParseFloat(corporate.FeeUser.AcceptPaymentCard, 64)
		if err != nil {
			return domain.Transaction{}, domain.Statement{}, utils.ErrorBadRequest(paramLog, utils.WrongAcceptCardFee, "Cannot convert accept payment card fee")
		}

		totalFee = int(float64(subAmount) * s)
	} else {
		s, err := strconv.ParseFloat(corporate.FeeCorporate.AcceptPaymentCard, 64)
		if err != nil {
			return domain.Transaction{}, domain.Statement{}, utils.ErrorBadRequest(paramLog, utils.WrongAcceptCardFee, "Cannot convert accept payment card fee")
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
		RequestId:        requestId,
	}

	statement := service.DepositTransactionStatement(
		balance.ID, transcation.Time, transcation.TransactionCode, subAmount)

	return transcation, statement, nil
}
