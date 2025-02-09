package transfer_bank

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
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

func (self UserTransferBank) Execute(paramLog *basic.ParamLog, corporate domain.Corporate, actor domain.ActorAble,
	to domain.TransactionObject, balanceID string, subAmount int, encryptedPIN string, notes string, externalID string, requestId string) (domain.Transaction, error) {

	balance, err := identifyBalance(paramLog, balanceID)
	if err != nil {
		return domain.Transaction{}, err
	}

	from, err := usecase.ActorObjectToActor(paramLog, balance.Owner.ToActorObject())
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

	transaction, transactionStatement := createTransaction(self.corporate, self.fromBalance, self.actor, self.from, to, subAmount, notes, externalID, requestId)

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(paramLog, corporate, self.fromBalance, transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	statements = append(statements, transactionStatement)
	statements = append(statements, feeStatement...)

	err = validateCurrency(paramLog, transaction, corporate)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = validationActor(paramLog, self.actor, self.fromBalance.ID.Hex(), self.pin)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = validationTransaction(paramLog, transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	self.transferBankBase.SetupGateway(&transaction)

	err = self.transactionUsecase.Commit(paramLog, statements, &transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = self.transferBankBase.CreateTransferGateway(paramLog, &transaction, requestId)
	go usecase.PublishTransferCallback(paramLog, corporate, transaction)
	return transaction, err
}

func identifyBalance(paramLog *basic.ParamLog, balanceID string) (domain.Balance, error) {
	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil {
		return domain.Balance{}, utils.ErrorBadRequest(paramLog, utils.InvalidBalanceID, "Balance id not found")
	}

	return balance, nil
}

func createTransaction(corporate domain.Corporate, balance domain.Balance, actor domain.ActorAble, from domain.TransactionObject,
	to domain.TransactionObject, subAmount int, notes string, externalID string, requestId string) (domain.Transaction, domain.Statement) {

	totalFee := 0
	if actor.GetActorType() == domain.ACTOR_TYPE_USER {
		totalFee = corporate.FeeUser.TransferBank
	} else {
		totalFee = corporate.FeeCorporate.TransferBank
	}

	transcation := domain.Transaction{
		TransactionCode:  utils.GenerateTransactionCode("1"),
		UserID:           actor.GetActorID(),
		CorporateID:      corporate.ID,
		Type:             domain.TRANSFER_BANK,
		Method:           domain.METHOD_BALANCE,
		FromBalanceID:    balance.ID,
		Actor:            actor.ToTransactionObject(),
		From:             from,
		To:               to,
		TotalFee:         totalFee,
		SubAmount:        subAmount,
		Amount:           subAmount + totalFee,
		Time:             time.Now().Format(os.Getenv("TIME_FORMAT")),
		Notes:            notes,
		Status:           domain.PENDING_STATUS,
		Unpaid:           false,
		ExternalID:       externalID,
		Currency:         corporate.Currency,
		RequestId:        requestId,
		GatewayReference: requestId,
	}

	statement := service.WithdrawTransactionStatement(
		balance.ID, transcation.Time, transcation.TransactionCode, subAmount)

	return transcation, statement
}

func validationActor(paramLog *basic.ParamLog, actor domain.ActorAble, balanceID string, pin string) error {

	err := usecase.ValidateActorPIN(paramLog, actor, pin)
	if err != nil {
		return err
	}

	err = usecase.ValidateAccessBalance(paramLog, actor, balanceID)
	if err != nil {
		return err
	}

	err = usecase.ValidateIsVerify(paramLog, actor)
	if err != nil {
		return err
	}

	return nil
}

func validationTransaction(paramLog *basic.ParamLog, transaction domain.Transaction) error {
	err := validateMaximum(paramLog, transaction)
	if err != nil {
		return err
	}

	err = validateMinimum(paramLog, transaction)
	if err != nil {
		return err
	}

	return nil
}
