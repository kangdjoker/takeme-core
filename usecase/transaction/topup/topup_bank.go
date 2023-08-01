package topup

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/gateway"
)

type TopupBank struct {
	corporate          domain.Corporate
	from               domain.Bank
	to                 domain.TransactionObject
	balance            domain.Balance
	amount             int
	currency           string
	reference          string
	transactionUsecase transaction.Base
}

func (self TopupBank) Execute(tag string, from domain.Bank, balanceID string, amount int,
	reference string, currency string, requestId string) (domain.Transaction, domain.Balance, error) {

	balance, owner, corporate, err := identifyBalance(balanceID)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	gateway := gateway.XenditGateway{}
	self.corporate = corporate
	self.from = from
	self.to = owner
	self.balance = balance
	self.amount = amount
	self.reference = reference
	self.transactionUsecase = transaction.Base{}
	self.currency = currency

	var statements []domain.Statement

	transaction, transactionStatement := createTransaction(self.corporate, self.balance, self.from,
		self.to, self.amount, self.reference, gateway, requestId)

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

	err = self.transactionUsecase.Commit(tag, statements, &transaction)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	go usecase.PublishTopupCallback(corporate, balance, transaction)

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

func createTransaction(corporate domain.Corporate, balance domain.Balance, from domain.Bank,
	to domain.TransactionObject, subAmount int, reference string, gateway gateway.Gateway, requestId string) (domain.Transaction, domain.Statement) {

	var totalFee = 0
	if balance.Owner.Type == domain.ACTOR_TYPE_USER {
		totalFee = corporate.FeeUser.Topup
	} else {
		totalFee = corporate.FeeCorporate.Topup
	}

	transcation := domain.Transaction{
		TransactionCode:  utils.GenerateTransactionCode("2"),
		CorporateID:      corporate.ID,
		Type:             domain.TOPUP,
		Method:           domain.METHOD_VA,
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
		ExternalID:       "",
		Gateway:          gateway.Name(),
		GatewayReference: reference,
		Currency:         corporate.Currency,
		RequestId:        requestId,
	}

	statement := service.DepositTransactionStatement(
		balance.ID, transcation.Time, transcation.TransactionCode, subAmount)

	return transcation, statement
}
