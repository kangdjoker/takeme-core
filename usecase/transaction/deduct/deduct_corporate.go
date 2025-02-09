package deduct

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

func (self DeductCorporate) Execute(paramLog *basic.ParamLog, corporate domain.Corporate, actor domain.ActorAble,
	toBalanceID string, fromBalanceID string, subAmount int, encryptedPIN string, externalID string, requestId string) (domain.Transaction, error) {

	fromBalance, _, _, err := identifyBalance(paramLog, fromBalanceID)
	if err != nil {
		return domain.Transaction{}, err
	}

	from, err := usecase.ActorObjectToActor(paramLog, fromBalance.Owner.ToActorObject())
	if err != nil {
		return domain.Transaction{}, err
	}

	toBalance, _, _, err := identifyBalance(paramLog, toBalanceID)
	if err != nil {
		return domain.Transaction{}, err
	}

	to, err := usecase.ActorObjectToActor(paramLog, toBalance.Owner.ToActorObject())
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
		self.toBalance, self.subAmount, self.externalID, requestId)

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(paramLog, corporate, self.fromBalance, transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	statements = append(statements, transactionStatement...)
	statements = append(statements, feeStatement...)

	err = validationActor(paramLog, self.actor, self.fromBalance, toBalance, self.pin)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = validateCurrency(paramLog, fromBalance, toBalance)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = validationTransaction(paramLog, transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = self.transactionUsecase.Commit(paramLog, statements, &transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	go usecase.PublishDeductCallback(paramLog, corporate, fromBalance, transaction)

	return transaction, nil
}

func createTransaction(corporate domain.Corporate, fromBalance domain.Balance, actor domain.ActorAble, from domain.TransactionObject,
	to domain.TransactionObject, toBalance domain.Balance, subAmount int, externalID string, requestId string) (domain.Transaction, []domain.Statement) {

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
		RequestId:       requestId,
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

func validationActor(paramLog *basic.ParamLog, actor domain.ActorAble, sourceBalance domain.Balance, targetBalance domain.Balance, pin string) error {

	err := usecase.ValidateActorPIN(paramLog, actor, pin)
	if err != nil {
		return err
	}

	err = usecase.ValidateAccessBalance(paramLog, actor, targetBalance.ID.Hex())
	if err != nil {
		return err
	}

	err = usecase.ValidateIsVerify(paramLog, actor)
	if err != nil {
		return err
	}

	err = validateDeductScope(paramLog, actor, sourceBalance)
	if err != nil {
		return err
	}

	return nil
}

func validateDeductScope(paramLog *basic.ParamLog, actor domain.ActorAble, balance domain.Balance) error {
	if balance.CorporateID.Hex() != actor.GetActorID().Hex() {
		return utils.ErrorBadRequest(paramLog, utils.InvalidDeductTarget, "Invalid deduct target")
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
