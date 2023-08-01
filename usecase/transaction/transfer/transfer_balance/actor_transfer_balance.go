package transfer_balance

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
)

type ActorTransferBalance struct {
	corporate          domain.Corporate
	actor              domain.ActorAble
	from               domain.TransactionObject
	to                 domain.TransactionObject
	fromBalance        domain.Balance
	toBalance          domain.Balance
	pin                string
	subAmount          int
	externalID         string
	transactionUsecase transaction.Base
	isTopuoType        bool
}

func (self ActorTransferBalance) Execute(tag string, corporate domain.Corporate, actor domain.ActorAble,
	toBalanceID string, fromBalanceID string, subAmount int, encryptedPIN string, externalID string, isTopupType bool, requestId string) (domain.Transaction, error) {

	fromBalance, err := identifyBalance(fromBalanceID)
	if err != nil {
		return domain.Transaction{}, err
	}

	from, err := usecase.ActorObjectToActor(fromBalance.Owner.ToActorObject())
	if err != nil {
		return domain.Transaction{}, err
	}

	toBalance, err := identifyBalance(toBalanceID)
	if err != nil {
		return domain.Transaction{}, err
	}

	to, err := usecase.ActorObjectToActor(toBalance.Owner.ToActorObject())
	if err != nil {
		return domain.Transaction{}, err
	}

	self.corporate = corporate
	self.actor = actor
	self.to = to.ToTransactionObject()
	self.from = from.ToTransactionObject()

	self.fromBalance = fromBalance
	self.toBalance = toBalance
	self.pin = encryptedPIN
	self.subAmount = subAmount
	self.externalID = externalID
	self.transactionUsecase = transaction.Base{}
	self.isTopuoType = isTopupType

	var statements []domain.Statement

	err = validateCurrency(fromBalance, toBalance)
	if err != nil {
		return domain.Transaction{}, err
	}

	transaction, transactionStatement := createTransaction(self.corporate, self.fromBalance, self.actor, self.from, self.to,
		self.toBalance, self.subAmount, self.externalID, isTopupType, requestId)

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(corporate, self.fromBalance, transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	statements = append(statements, transactionStatement...)
	statements = append(statements, feeStatement...)

	if self.actor.GetActorType() == domain.ACTOR_TYPE_USER {
		err = validationActorUser(self.actor, self.fromBalance.ID.Hex(), self.pin)
		if err != nil {
			return domain.Transaction{}, err
		}
	} else {
		err = validationActorCorporate(self.actor, fromBalance, corporate, self.pin)
		if err != nil {
			return domain.Transaction{}, err
		}
	}

	err = validationTransaction(transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = self.transactionUsecase.Commit(tag, statements, &transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	go usecase.PublishTopupCallback(corporate, toBalance, transaction)

	return transaction, nil
}

func createTransaction(corporate domain.Corporate, fromBalance domain.Balance, actor domain.ActorAble, from domain.TransactionObject,
	to domain.TransactionObject, toBalance domain.Balance, subAmount int, externalID string, isTopupType bool, requestId string) (domain.Transaction, []domain.Statement) {

	totalFee := 0
	if actor.GetActorType() == domain.ACTOR_TYPE_USER {
		totalFee = corporate.FeeUser.TransferBalance
	} else {
		totalFee = corporate.FeeCorporate.TransferBalance
	}

	transactionType := domain.TRANSFER_WALLET
	if isTopupType == true {
		transactionType = domain.TOPUP
	}

	transaction := domain.Transaction{
		TransactionCode: utils.GenerateTransactionCode("1"),
		UserID:          actor.GetActorID(),
		CorporateID:     corporate.ID,
		Type:            transactionType,
		Method:          domain.METHOD_BALANCE,
		FromBalanceID:   fromBalance.ID,
		ToBalanceID:     toBalance.ID,
		Actor:           actor.ToTransactionObject(),
		From:            from,
		To:              to,
		TotalFee:        totalFee,
		SubAmount:       subAmount,
		Amount:          subAmount + totalFee,
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
		Notes:           "",
		Status:          domain.COMPLETED_STATUS,
		Unpaid:          false,
		ExternalID:      externalID,
		Currency:        corporate.Currency,
		RequestId:       requestId,
	}

	var statements []domain.Statement

	fromStatement := service.WithdrawTransactionStatement(
		fromBalance.ID, transaction.Time, transaction.TransactionCode, subAmount)

	toStatement := service.DepositTransactionStatement(
		toBalance.ID, transaction.Time, transaction.TransactionCode, subAmount)

	statements = append(statements, fromStatement)
	statements = append(statements, toStatement)

	return transaction, statements
}

func identifyBalance(balanceID string) (domain.Balance, error) {
	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil {
		return domain.Balance{}, utils.ErrorBadRequest(utils.InvalidBalanceID, "Balance id not found")
	}

	return balance, nil
}

func validationActorUser(actor domain.ActorAble, balanceID string, pin string) error {

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

func validationActorCorporate(actor domain.ActorAble, balance domain.Balance, corporate domain.Corporate, pin string) error {

	err := usecase.ValidateActorPIN(actor, pin)
	if err != nil {
		return err
	}

	if balance.CorporateID != balance.CorporateID {
		return utils.ErrorBadRequest(utils.InvalidBalanceAccess, "Invalid balance access")
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
