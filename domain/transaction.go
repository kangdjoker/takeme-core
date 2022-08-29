package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TRANSACTION_COLLECTION string = "transaction"

const (
	REFUND_STATUS    = "Refund"
	COMPLETED_STATUS = "Completed"
	PENDING_STATUS   = "Pending"
	FAILED_STATUS    = "Failed"
)

const (
	TOPUP               = "TOPUP"
	DEDUCT              = "DEDUCT"
	ACCEPT_PAYMENT_CARD = "ACCEPT_PAYMENT_CARD"
	TRANSFER_WALLET     = "TRANSFER_TO_WALLET"
	TRANSFER_CASH       = "TRANSFER_TO_CASH"
	TRANSFER_BANK       = "TRANSFER_TO_BANK"
	PAY_QR              = "PAY_QR"
	BILLER              = "PAY_BILLER"
)

const (
	METHOD_VA      = "Virtual Account"
	METHOD_CARD    = "Credit Card"
	METHOD_BALANCE = "Wallet Balance"
)

const (
	BANK_OBJECT      = "BANK_ACCOUNT"
	CARD_OBJECT      = "CARD_ACCOUNT"
	WALLET_OBJECT    = "WALLET_ACCOUNT"
	BILLER_OBJECT    = "BILLER"
	CORPORATE_OBJECT = "CORPORATE_ACCOUNT"
	PERSON_OBJECT    = "PERSON"
)

type TransactionObject struct {
	Type            string `json:"type" bson:"type,omitempty"`
	InstitutionCode string `json:"institution_code" bson:"institution_code,omitempty"`
	Name            string `json:"name" bson:"name,omitempty"`
	AccountNumber   string `json:"account_number" bson:"account_number,omitempty"`
}

// transaction able interface
func (model TransactionObject) GetType() string {
	return model.Type
}

func (model TransactionObject) GetInstitutionCode() string {
	return model.InstitutionCode
}

func (model TransactionObject) GetName() string {
	return model.Name
}

func (model TransactionObject) GetAccountNumber() string {
	return model.AccountNumber
}

type TransactionAble interface {
	GetType() string
	GetInstitutionCode() string
	GetName() string
	GetAccountNumber() string
	ToTransactionObject() TransactionObject
}

type Transaction struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TransactionCode   string             `json:"transaction_code" bson:"transaction_code,omitempty"`
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	CorporateID       primitive.ObjectID `json:"corporate_id" bson:"corporate_id,omitempty"`
	Type              string             `json:"type" bson:"type,omitempty"`
	Method            string             `json:"method" bson:"method,omitempty"`
	FromBalanceID     primitive.ObjectID `json:"from_balance_id" bson:"from_balance_id,omitempty"`
	ToBalanceID       primitive.ObjectID `json:"to_balance_id" bson:"to_balance_id,omitempty"`
	Actor             TransactionObject  `json:"actor" bson:"actor,omitempty"`
	From              TransactionObject  `json:"from" bson:"from,omitempty"`
	To                TransactionObject  `json:"to" bson:"to,omitempty"`
	DetailsFee        []DetailFee        `json:"details_fee" bson:"details_fee"`
	TotalFee          int                `json:"total_fee" bson:"total_fee"`
	SubAmount         int                `json:"sub_amount" bson:"sub_amount"`
	Amount            int                `json:"amount" bson:"amount"`
	Time              string             `json:"time" bson:"time,omitempty"`
	Notes             string             `json:"notes" bson:"notes"`
	Status            string             `json:"status" bson:"status"`
	RerunStatus       string             `json:"rerun_status" bson:"rerun_status"`
	Unpaid            bool               `json:"unpaid" bson:"unpaid"`
	Description       string             `json:"description" bson:"description"`
	CashoutCode       string             `json:"cashout_code" bson:"cashout_code,omitempty"`
	ExternalID        string             `json:"external_id" bson:"external_id"`
	Gateway           string             `json:"gateway" bson:"gateway"`
	GatewayReference  string             `json:"gateway_reference" bson:"gateway_reference"`
	GatewayStrategies []GatewayStrategy  `json:"gateway_strategies" bson:"gateway_strategies"`
	GatewayHistories  []GatewayHistory   `json:"gateway_histories" bson:"gateway_histories"`
	Currency          string             `json:"currency" bson:"currency,omitempty"`
}

// Interface for mongo document result
func (domain *Transaction) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *Transaction) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (model *Transaction) CollectionName() string {
	return TRANSACTION_COLLECTION
}

type DetailFee struct {
	CorporateID primitive.ObjectID `json:"corporate_id" bson:"corporate_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Amount      int                `json:"amount" bson:"amount"`
}

type GatewayHistory struct {
	Code      string `json:"code" bson:"code"`
	Reference string `json:"reference" bson:"reference"`
	Time      string `json:"time" bson:"time,omitempty"`
}

type GatewayStrategy struct {
	Code       string `json:"code" bson:"code"`
	IsExecuted bool   `json:"is_executed" bson:"is_executed"`
}
