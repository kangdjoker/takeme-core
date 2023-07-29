package biller

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/domain/dto"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
)

type BPJSTKBiller struct {
	corporate   domain.Corporate
	actor       domain.ActorAble
	to          domain.TransactionObject
	fromBalance domain.Balance
	// toBalance          domain.Balance
	pin                string
	externalID         string
	transactionUsecase transaction.Base
	billerBase         BillerBase
	paymentCode        string
	currency           string
}

func (biller BPJSTKBiller) Inquiry(paymentCode string, currency string) (FusBPJSInqResponse, error) {
	return biller.billerBase.BillerInquiryBPJSTKPMI(paymentCode, currency)
}
func (self BPJSTKBiller) Execute(corporate domain.Corporate, actor domain.ActorAble,
	to domain.TransactionObject, balanceID string, encryptedPIN string, externalID string,
	paymentCode string, currency string) (domain.Transaction, interface{}, error) {

	balance, err := identifyBalance(balanceID)
	if err != nil {
		return domain.Transaction{}, nil, err
	}

	self.corporate = corporate
	self.actor = actor
	self.to = to
	self.pin = encryptedPIN
	self.externalID = externalID
	self.fromBalance = balance
	self.transactionUsecase = transaction.Base{}
	self.billerBase = BillerBase{}

	var statements []domain.Statement

	transaction, transactionStatement := createTransaction(self.corporate, self.fromBalance, self.actor, to, 80000, externalID)

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(corporate, self.fromBalance, transaction)
	if err != nil {
		return domain.Transaction{}, nil, err
	}

	statements = append(statements, transactionStatement)
	statements = append(statements, feeStatement...)

	err = validationActor(self.actor, self.fromBalance.ID.Hex(), self.pin)
	if err != nil {
		return domain.Transaction{}, nil, err
	}

	err, ref := self.billerBase.BillerPayBPJSTKPMI(transaction, paymentCode, currency)
	if err != nil {
		return domain.Transaction{}, nil, err
	}

	transaction.GatewayReference = ref

	err = self.transactionUsecase.Commit(statements, &transaction)
	if err != nil {
		return domain.Transaction{}, nil, err
	}

	return transaction, dto.BPJSTKPMI{
		Name:              "Laura Basuki",
		KPJNumber:         "4837824782378",
		DateOfBirth:       "02-05-1984",
		PaymentCode:       paymentCode,
		MonthOfProtection: "14-07-2022 s/d 14-08-2022",
		Reference:         "2056080",
		JKK:               "30000",
		JKM:               "50000",
		JHT:               "0",
		SubAmount:         "80000",
		TotalFee:          "3500",
		Amount:            "83500",
		CurrencyCode:      currency,
		FixedRate:         "0.00031",
		LocalInvoice:      "24.8",
	}, err
}

func identifyBalance(balanceID string) (domain.Balance, error) {
	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil {
		return domain.Balance{}, utils.ErrorBadRequest(utils.InvalidBalanceID, "Balance id not found")
	}

	return balance, nil
}

func createTransaction(corporate domain.Corporate, balance domain.Balance, from domain.ActorAble,
	to domain.TransactionObject, subAmount int, externalID string) (domain.Transaction, domain.Statement) {

	totalFee := 0
	if from.GetActorType() == domain.ACTOR_TYPE_USER {
		totalFee = corporate.FeeUser.TransferBank
	} else {
		totalFee = corporate.FeeCorporate.TransferBank
	}

	transcation := domain.Transaction{
		TransactionCode: utils.GenerateTransactionCode("2"),
		UserID:          from.GetActorID(),
		CorporateID:     corporate.ID,
		Type:            domain.BILLER,
		Method:          domain.METHOD_BALANCE,
		FromBalanceID:   balance.ID,
		From:            from.ToTransactionObject(),
		To:              to,
		TotalFee:        totalFee,
		SubAmount:       subAmount,
		Amount:          subAmount + totalFee,
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
		Notes:           "",
		Status:          domain.COMPLETED_STATUS,
		Unpaid:          false,
		ExternalID:      externalID,
		Currency:        "idr",
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
