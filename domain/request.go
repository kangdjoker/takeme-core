package domain

import (
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	REQUEST_STATUS_PENDING   = "PENDING"
	REQUEST_STATUS_COMPLETED = "COMPLETED"
	REQUEST_STATUS_REJECTED  = "REJECTED"
)

const REQUEST_COLLECTION string = "request"

type Request struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	CorporateID     primitive.ObjectID `json:"corporate_id" bson:"corporate_id,omitempty"`
	Status          string             `json:"status" bson:"status,omitempty"`
	FromBalanceID   primitive.ObjectID `json:"from_balance_id" bson:"from_balance_id,omitempty"`
	ToBalanceID     primitive.ObjectID `json:"to_balance_id" bson:"to_balance_id,omitempty"`
	From            TransactionObject  `json:"from" bson:"from,omitempty"`
	To              TransactionObject  `json:"to" bson:"to,omitempty"`
	Amount          int                `json:"amount" bson:"amount"`
	TransactionCode string             `json:"transaction_code" bson:"transaction_code,omitempty"`
	Time            string             `json:"time" bson:"time,omitempty"`
	IsRead          bool               `json:"is_read" bson:"is_read"`
	Message         string             `json:"message" bson:"message,omitempty"`
}

func CreateRequest(corporateID primitive.ObjectID, fromUser User,
	toUser User, amount int) (Request, error) {

	return Request{
		UserID:      toUser.ID,
		CorporateID: corporateID,
		Status:      REQUEST_STATUS_PENDING,
		From: TransactionObject{
			Type:            WALLET_OBJECT,
			InstitutionCode: fromUser.ID.Hex(),
			Name:            fromUser.FullName,
			AccountNumber:   fromUser.PhoneNumber,
		},
		To: TransactionObject{
			Type:            WALLET_OBJECT,
			InstitutionCode: toUser.ID.Hex(),
			Name:            toUser.FullName,
			AccountNumber:   toUser.PhoneNumber,
		},
		Amount: amount,
		Time:   time.Now().Format(os.Getenv("TIME_FORMAT")),
		IsRead: false,
	}, nil
}

func (domain *Request) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *Request) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (domain *Request) CollectionName() string {
	return REQUEST_COLLECTION
}
