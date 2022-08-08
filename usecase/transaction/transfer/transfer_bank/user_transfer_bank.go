package transfer_bank

import (
	"os"
	"time"

	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/service"
	"github.com/takeme-id/core/usecase"
	"github.com/takeme-id/core/usecase/transaction"
	"github.com/takeme-id/core/utils"
)

type UserTransferBank struct {
	corporate   domain.Corporate
	actor       domain.ActorAble
	from        domain.TransactionObject
	to          domain.TransactionObject
	fromBalance domain.Balance
	// toBalance          domain.Balance
	pin                string
	subAmount          int
	externalID         string
	transactionUsecase transaction.Base
	transferBankBase   TransferBank
}

func (self UserTransferBank) Execute(corporate domain.Corporate, actor domain.ActorAble,
	to domain.TransactionObject, balanceID string, subAmount int, encryptedPIN string, externalID string) (domain.Transaction, error) {

	balance, err := identifyBalance(balanceID)
	if err != nil {
		return domain.Transaction{}, err
	}

	from, err := usecase.ActorObjectToActor(balance.Owner.ToActorObject())
	if err != nil {
		return domain.Transaction{}, err
	}

	self.corporate = corporate
	self.actor = actor
	self.to = to
	self.subAmount = subAmount
	self.pin = encryptedPIN
	self.externalID = externalID
	self.from = from.ToTransactionObject()
	self.fromBalance = balance
	self.transactionUsecase = transaction.Base{}
	self.transferBankBase = TransferBank{}

	var statements []domain.Statement

	transaction, transactionStatement := createTransaction(self.corporate, self.fromBalance, self.actor, self.from, to, subAmount, externalID)

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(corporate, self.fromBalance, transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	statements = append(statements, transactionStatement)
	statements = append(statements, feeStatement...)

	err = validationActor(self.actor, self.fromBalance.ID.Hex(), self.pin)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = validationTransaction(transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	self.transferBankBase.SetupGateway(&transaction)

	err = self.transactionUsecase.Commit(statements, &transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	go self.transferBankBase.CreateTransferGateway(transaction)

	return transaction, nil
}

func identifyBalance(balanceID string) (domain.Balance, error) {
	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil {
		return domain.Balance{}, utils.ErrorBadRequest(utils.InvalidBalanceID, "Balance id not found")
	}

	return balance, nil
}

func createTransaction(corporate domain.Corporate, balance domain.Balance, actor domain.ActorAble, from domain.TransactionObject,
	to domain.TransactionObject, subAmount int, externalID string) (domain.Transaction, domain.Statement) {

	totalFee := 0
	if actor.GetActorType() == domain.ACTOR_TYPE_USER {
		totalFee = corporate.FeeUser.TransferBank
	} else {
		totalFee = corporate.FeeCorporate.TransferBank
	}

	transcation := domain.Transaction{
		TransactionCode: utils.GenerateTransactionCode("1"),
		UserID:          actor.GetActorID(),
		CorporateID:     corporate.ID,
		Type:            domain.TRANSFER_BANK,
		Method:          domain.METHOD_BALANCE,
		FromBalanceID:   balance.ID,
		Actor:           actor.ToTransactionObject(),
		From:            from,
		To:              to,
		TotalFee:        totalFee,
		SubAmount:       subAmount,
		Amount:          subAmount + totalFee,
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
		Notes:           "",
		Status:          domain.PENDING_STATUS,
		Unpaid:          false,
		ExternalID:      externalID,
	}

	statement := service.WithdrawTransactionStatement(
		balance.ID, transcation.Time, transcation.TransactionCode, subAmount)

	return transcation, statement
}

func validationActor(actor domain.ActorAble, balanceID string, pin string) error {

	err := usecase.ValidateActorPIN(actor, pin)
	if err != nil {
		return err
	}

	err = usecase.ValidateAccessBalance(actor, balanceID)
	if err != nil {
		return err
	}

	err = usecase.ValidateIsVerify(actor)
	if err != nil {
		return err
	}

	return nil
}

func validationTransaction(transaction domain.Transaction) error {
	err := validateMaximum(transaction)
	if err != nil {
		return err
	}

	err = validateMinimum(transaction)
	if err != nil {
		return err
	}

	return nil
}
