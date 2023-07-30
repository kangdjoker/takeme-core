package topup

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
)

type TopupManual struct {
	corporate          domain.Corporate
	to                 domain.TransactionObject
	balance            domain.Balance
	amount             int
	currency           string
	remark             string
	transactionUsecase transaction.Base
}

func (tm TopupManual) Execute(balanceID string, amount int,
	remark string, currency string) (domain.Transaction, domain.Balance, error) {

	balance, owner, corporate, err := identifyBalance(balanceID)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	tm.corporate = corporate
	tm.to = owner
	tm.balance = balance
	tm.amount = amount
	tm.remark = remark
	tm.transactionUsecase = transaction.Base{}
	tm.currency = currency

	var statements []domain.Statement

	transaction, transactionStatement := createTransactionTopupManual(tm.corporate, tm.balance,
		tm.to, tm.amount, tm.remark)

	statements = append(statements, transactionStatement)

	err = validateCurrency(tm.currency, balance)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	err = tm.transactionUsecase.Commit(statements, &transaction)
	if err != nil {
		return domain.Transaction{}, domain.Balance{}, err
	}

	return transaction, balance, nil
}

func createTransactionTopupManual(corporate domain.Corporate, balance domain.Balance,
	to domain.TransactionObject, subAmount int, remark string) (domain.Transaction, domain.Statement) {

	var totalFee = 0
	newUuid := uuid.New().String()
	transcation := domain.Transaction{
		TransactionCode:  utils.GenerateTransactionCode("2"),
		CorporateID:      corporate.ID,
		Type:             domain.TOPUP,
		Method:           domain.METHOD_VA,
		ToBalanceID:      balance.ID,
		FromBalanceID:    balance.ID,
		From:             domain.TransactionObject{},
		To:               to,
		TotalFee:         totalFee,
		SubAmount:        subAmount,
		Amount:           subAmount - totalFee,
		Time:             time.Now().Format(os.Getenv("TIME_FORMAT")),
		Notes:            remark,
		Status:           domain.COMPLETED_STATUS,
		Unpaid:           false,
		ExternalID:       "",
		Gateway:          "Manual",
		GatewayReference: newUuid,
		Currency:         corporate.Currency,
		RequestId:        newUuid,
	}

	statement := service.DepositTransactionStatement(
		balance.ID, transcation.Time, transcation.TransactionCode, subAmount)

	return transcation, statement
}
