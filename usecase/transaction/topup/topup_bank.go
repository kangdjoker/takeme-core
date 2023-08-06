package topup

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
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

func (self TopupBank) Execute(paramLog *basic.ParamLog, from domain.Bank, balanceID string, amount int,
	reference string, currency string, requestId string) (domain.Transaction, domain.Balance, error) {

	balance, owner, corporate, err := identifyBalance(paramLog, balanceID)
	if err != nil {
		basic.LogError2(paramLog, "identifyBalance", err)
		return domain.Transaction{}, domain.Balance{}, err
	}
	basic.LogInformation2(paramLog, "balance", balance)
	basic.LogInformation2(paramLog, "owner", owner)
	basic.LogInformation2(paramLog, "corporate", corporate)

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
	basic.LogInformation2(paramLog, "transaction", transaction)
	basic.LogInformation2(paramLog, "transactionStatement", transactionStatement)

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(paramLog, corporate, balance, transaction)
	if err != nil {
		basic.LogError2(paramLog, "CreateFeeStatement", err)
		return domain.Transaction{}, domain.Balance{}, err
	}
	basic.LogInformation2(paramLog, "feeStatement", feeStatement)

	statements = append(statements, transactionStatement)
	statements = append(statements, feeStatement...)

	err = validateCurrency(paramLog, balance, balance)
	if err != nil {
		basic.LogError2(paramLog, "validateCurrency", err)
		return domain.Transaction{}, domain.Balance{}, err
	}
	basic.LogInformation(paramLog, "validateCurrency.Success")

	err = self.transactionUsecase.Commit(paramLog, statements, &transaction)
	if err != nil {
		basic.LogError2(paramLog, "transactionUsecase.Commit", err)
		return domain.Transaction{}, domain.Balance{}, err
	}
	basic.LogInformation(paramLog, "transactionUsecase.Commit.Success")

	go usecase.PublishTopupCallback(paramLog, corporate, balance, transaction)

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
