package deduct

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
)

type DeductCorporate struct {
	corporate          domain.Corporate
	actor              domain.ActorAble
	to                 domain.TransactionObject
	from               domain.TransactionObject
	fromBalance        domain.Balance
	toBalance          domain.Balance
	pin                string
	subAmount          int
	externalID         string
	transactionUsecase transaction.Base
}

func (self DeductCorporate) Execute(corporate domain.Corporate, actor domain.ActorAble,
	toBalanceID string, fromBalanceID string, subAmount int, encryptedPIN string, externalID string) (domain.Transaction, error) {

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

	var statements []domain.Statement

	transaction, transactionStatement := createTransaction(self.corporate, self.fromBalance, self.actor, self.from, self.to,
		self.toBalance, self.subAmount, self.externalID)

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(corporate, self.fromBalance, transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	statements = append(statements, transactionStatement...)
	statements = append(statements, feeStatement...)

	err = validationActor(self.actor, self.fromBalance, toBalance, self.pin)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = validateCurrency(fromBalance, toBalance)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = validationTransaction(transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = self.transactionUsecase.Commit(statements, &transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	go usecase.PublishDeductCallback(corporate, fromBalance, transaction)

	return transaction, nil
}

func createTransaction(corporate domain.Corporate, fromBalance domain.Balance, actor domain.ActorAble, from domain.TransactionObject,
	to domain.TransactionObject, toBalance domain.Balance, subAmount int, externalID string) (domain.Transaction, []domain.Statement) {

	totalFee := 0
	if actor.GetActorType() == domain.ACTOR_TYPE_USER {
		totalFee = corporate.FeeUser.Deduct
	} else {
		totalFee = corporate.FeeCorporate.Deduct
	}

	transaction := domain.Transaction{
		TransactionCode: utils.GenerateTransactionCode("1"),
		UserID:          actor.GetActorID(),
		CorporateID:     corporate.ID,
		Type:            domain.DEDUCT,
		Method:          domain.METHOD_BALANCE,
		FromBalanceID:   toBalance.ID,
		ToBalanceID:     fromBalance.ID,
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
	}

	var statements []domain.Statement

	fromStatement := service.DepositTransactionStatement(
		toBalance.ID, transaction.Time, transaction.TransactionCode, subAmount)

	toStatement := service.WithdrawTransactionStatement(
		fromBalance.ID, transaction.Time, transaction.TransactionCode, subAmount)

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

func validationActor(actor domain.ActorAble, sourceBalance domain.Balance, targetBalance domain.Balance, pin string) error {

	err := usecase.ValidateActorPIN(actor, pin)
	if err != nil {
		return err
	}

	err = usecase.ValidateAccessBalance(actor, targetBalance.ID.Hex())
	if err != nil {
		return err
	}

	err = usecase.ValidateIsVerify(actor)
	if err != nil {
		return err
	}

	err = validateDeductScope(actor, sourceBalance)
	if err != nil {
		return err
	}

	return nil
}

func validateDeductScope(actor domain.ActorAble, balance domain.Balance) error {
	if balance.CorporateID.Hex() != actor.GetActorID().Hex() {
		return utils.ErrorBadRequest(utils.InvalidDeductTarget, "Invalid deduct target")
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
