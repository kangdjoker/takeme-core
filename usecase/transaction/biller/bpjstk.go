package biller

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/domain/dto"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
	"github.com/kangdjoker/takeme-core/usecase/transaction"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
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

func (biller BPJSTKBiller) Inquiry(paramLog *basic.ParamLog, paymentCode string, currency string, requestId string) (FusBPJSInqResponse, error) {
	return biller.billerBase.BillerInquiryBPJSTKPMI(paramLog, paymentCode, currency, requestId)
}
func (self BPJSTKBiller) Execute(paramLog *basic.ParamLog, corporate domain.Corporate, actor domain.ActorAble,
	to domain.TransactionObject, balanceID string, encryptedPIN string, externalID string,
	paymentCode string, currency string, requestId string) (domain.Transaction, interface{}, error) {

	balance, err := identifyBalance(paramLog, balanceID)
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

	//Ambil data Inqiry dulu
	resInquiry, err := self.billerBase.BillerInquiryBPJSTKPMI(paramLog, paymentCode, currency, strings.ReplaceAll(uuid.New().String(), "-", ""))
	if err != nil {
		return domain.Transaction{}, nil, err
	}
	if resInquiry.Status == "1" {
		return domain.Transaction{}, nil, errors.New(resInquiry.Data1)
	} else if resInquiry.Status != "0" {
		return domain.Transaction{}, nil, errors.New("unknown error")
	}

	to.Name = resInquiry.Data2
	totalBayar, err := strconv.Atoi(resInquiry.TotalBayarRupiah)
	if err != nil {
		return domain.Transaction{}, nil, err
	}
	if totalBayar == 0 {
		return domain.Transaction{}, nil, errors.New("tidak mendapatkan data inquiry")
	}

	transaction, transactionStatement := createTransaction(paramLog, self.corporate, self.fromBalance, self.actor, to, totalBayar, externalID, requestId)

	feeStatement, err := self.transactionUsecase.CreateFeeStatement(paramLog, corporate, self.fromBalance, transaction)
	if err != nil {
		return domain.Transaction{}, nil, err
	}

	statements = append(statements, transactionStatement)
	statements = append(statements, feeStatement...)

	err = validationActor(paramLog, self.actor, self.fromBalance.ID.Hex(), self.pin)
	if err != nil {
		return domain.Transaction{}, nil, err
	}

	err = self.transactionUsecase.Commit(paramLog, statements, &transaction)
	if err != nil {
		basic.LogInformation(paramLog, "Error: "+err.Error())
		return domain.Transaction{}, nil, err
	}

	resPayment, err := self.billerBase.BillerPayBPJSTKPMI(paramLog, transaction, paymentCode, currency, requestId)
	if err != nil {
		return domain.Transaction{}, nil, err
	}
	if resPayment.Status == "1" {
		return domain.Transaction{}, nil, utils.ErrorUnprocessableEntity(paramLog, 0, resPayment.Data1)
	} else if resPayment.Status != "0" {
		return domain.Transaction{}, nil, errors.New("unknown error")
	}

	//UPDATE Reff
	transaction.GatewayReference = resPayment.Ftrxid
	self.transactionUsecase.UpdatingTransactionDetail(paramLog, &transaction)

	return transaction, dto.BPJSTKPMI{
		Name:              resPayment.Data2,
		KPJNumber:         resPayment.Data1,
		DateOfBirth:       resInquiry.Data3,
		PaymentCode:       paymentCode,
		MonthOfProtection: resPayment.Blth,
		Reference:         resPayment.Reff,
		JKK:               resInquiry.Jkk,
		JKM:               resInquiry.Jkm,
		JHT:               resInquiry.Jht,
		SubAmount:         resInquiry.Tagihan,
		TotalFee:          resInquiry.Admin,
		Amount:            resInquiry.TotalBayarRupiah,
		CurrencyCode:      currency,
		FixedRate:         "0.0",
		LocalInvoice:      "0.0",
	}, err
}

func identifyBalance(paramLog *basic.ParamLog, balanceID string) (domain.Balance, error) {
	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil {
		return domain.Balance{}, utils.ErrorBadRequest(paramLog, utils.InvalidBalanceID, "Balance id not found")
	}

	return balance, nil
}

func createTransaction(paramLog *basic.ParamLog, corporate domain.Corporate, balance domain.Balance, from domain.ActorAble,
	to domain.TransactionObject, subAmount int, externalID string, requestId string) (domain.Transaction, domain.Statement) {

	totalFee := 0
	if from.GetActorType() == domain.ACTOR_TYPE_USER {
		totalFee = corporate.FeeUser.TransferBank
	} else {
		totalFee = corporate.FeeCorporate.TransferBank
	}

	transcation := domain.Transaction{
		TransactionCode:  utils.GenerateTransactionCode("2"),
		UserID:           from.GetActorID(),
		CorporateID:      corporate.ID,
		Type:             domain.BILLER,
		Method:           domain.METHOD_BALANCE,
		FromBalanceID:    balance.ID,
		From:             from.ToTransactionObject(),
		To:               to,
		TotalFee:         totalFee,
		SubAmount:        subAmount,
		Amount:           subAmount + totalFee,
		Time:             time.Now().Format(os.Getenv("TIME_FORMAT")),
		Notes:            "",
		Status:           domain.COMPLETED_STATUS,
		Unpaid:           false,
		ExternalID:       externalID,
		Currency:         "idr",
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
