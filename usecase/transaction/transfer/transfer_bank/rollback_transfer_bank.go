package transfer_bank

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
)

type RollbackTransferBank struct {
	corporate          domain.Corporate
	balance            domain.Balance
	transaction        domain.Transaction
	transactionUsecase transaction.Base
}

func (self *RollbackTransferBank) Initialize(rollbackTransaction domain.Transaction) error {
	corporate, err := service.CorporateByIDNoSession(rollbackTransaction.CorporateID.Hex())
	if err != nil {
		return err
	}

	self.corporate = corporate

	balance, err := service.BalanceByIDNoSession(rollbackTransaction.FromBalanceID.Hex())
	if err != nil {
		return err
	}

	self.balance = balance

	self.transaction = rollbackTransaction
	self.transactionUsecase = transaction.Base{}

	return nil
}

func (self *RollbackTransferBank) ExecuteRollback() error {
	transactionStatement := service.DepositTransactionStatement(
		self.balance.ID, time.Now().Format(os.Getenv("TIME_FORMAT")),
		self.transaction.TransactionCode,
		self.transaction.SubAmount)

	feeStatements, err := self.transactionUsecase.RollbackFeeStatement(self.corporate, self.balance, self.transaction)
	if err != nil {
		return err
	}

	var statements []domain.Statement
	statements = append(statements, transactionStatement)
	statements = append(statements, feeStatements...)

	err = self.transactionUsecase.CommitRollback(statements)
	if err != nil {
		return err
	}

	return nil
}
