package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const BALANCE_COLLECTION string = "balance"

const (
	ACCESS_BALANCE_VIEW_ONLY = "View Only"
	ACCESS_BALANCE_SHARED    = "Shared"
	ACCESS_BALANCE_OWNER     = "Owner"
)

type Balance struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CorporateID primitive.ObjectID `json:"corporate_id" bson:"corporate_id,omitempty"`
	Owner       ActorObject        `json:"owner" bson:"owner,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Amount      int                `json:"amount" bson:"amount"`
	VA          []VirtualAccount   `json:"va" bson:"va,omitempty"`
}

type VirtualAccount struct {
	BankCode      string `json:"bank_code" bson:"bank_code,omitempty"`
	AccountNumber string `json:"account_number" bson:"account_number,omitempty"`
}

func (domain *Balance) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *Balance) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (domain *Balance) CollectionName() string {
	return BALANCE_COLLECTION
}
