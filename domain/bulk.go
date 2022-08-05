package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const BULK_TRANSFER_COLLECTION string = "bulk_transfer"
const BULK_INQUIRY_COLLECTION string = "bulk_inquiry"

const (
	BULK_UNEXECUTED_STATUS = "Unexecuted"
	BULK_PROGRESS_STATUS   = "Progress"
	BULK_COMPLETED_STATUS  = "Completed"
)

type BulkTransfer struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CorporateID  primitive.ObjectID `json:"corporate_id" bson:"corporate_id,omitempty"`
	BalanceID    primitive.ObjectID `json:"balance_id" bson:"balance_id,omitempty"`
	Reference    string             `json:"reference" bson:"reference,omitempty"`
	Owner        ActorObject        `json:"owner" bson:"owner,omitempty"`
	Time         string             `json:"time" bson:"time,omitempty"`
	SubAmount    int                `json:"sub_amount" bson:"sub_amount,omitempty"`
	Amount       int                `json:"amount" bson:"amount,omitempty"`
	List         []Transfer         `json:"list" bson:"list,omitempty"`
	TotalList    int                `json:"total_list" bson:"total_list,omitempty"`
	Status       string             `json:"status" bson:"status,omitempty"`
	FailedNumber []int              `json:"failedNumber" bson:"failedNumber,omitempty"`
}

type Transfer struct {
	Number          int    `json:"number" bson:"number,omitempty"`
	ToBankAccount   Bank   `json:"to_bank_account" bson:"to_bank_account,omitempty"`
	Type            string `json:"type" bson:"type,omitempty"`
	Notes           string `json:"notes" bson:"notes,omitempty"`
	Amount          int    `json:"amount" bson:"amount,omitempty"`
	ExternalID      string `json:"external_id" bson:"external_id,omitempty"`
	TransactionCode string `json:"transaction_code" bson:"transaction_code"`
	Reason          string `json:"reason" bson:"reason,omitempty"`
}

type BulkInquiry struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CorporateID primitive.ObjectID `json:"corporate_id" bson:"corporate_id,omitempty"`
	Reference   string             `json:"reference" bson:"reference,omitempty"`
	Owner       ActorObject        `json:"owner" bson:"owner,omitempty"`
	Time        string             `json:"time" bson:"time,omitempty"`
	List        []Inquiry          `json:"list" bson:"list,omitempty"`
	TotalList   int                `json:"total_list" bson:"total_list,omitempty"`
	Status      string             `json:"status" bson:"status,omitempty"`
}

type Inquiry struct {
	Number        int    `json:"number" bson:"number,omitempty"`
	AccountName   string `json:"account_name" bson:"account_name"`
	AccountNumber string `json:"account_number" bson:"account_number"`
	BankName      string `json:"bank_name" bson:"bank_name"`
	Valid         bool   `json:"valid"  bson:"valid"`
	Reason        string `json:"reason"  bson:"reason"`
}

// Interface for mongo document result
func (domain *BulkTransfer) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *BulkTransfer) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (model *BulkTransfer) CollectionName() string {
	return BULK_TRANSFER_COLLECTION
}

// Interface for mongo document result
func (domain *BulkInquiry) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *BulkInquiry) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (model *BulkInquiry) CollectionName() string {
	return BULK_INQUIRY_COLLECTION
}
