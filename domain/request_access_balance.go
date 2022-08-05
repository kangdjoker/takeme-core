package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	REQUEST_ACCESS_BALANCE_STATUS_PENDING = "Pending"
	REQUEST_ACCESS_BALANCE_STATUS_APPROVE = "Approve"
	REQUEST_ACCESS_BALANCE_STATUS_REJECT  = "Reject"
)

const RAB_COLLECTION_NAME string = "request_access_balance"

type RequestAccessBalance struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CorporateID      primitive.ObjectID `json:"corporate_id" bson:"corporate_id,omitempty"`
	BalanceID        primitive.ObjectID `json:"balance_id" bson:"balance_id,omitempty"`
	Time             string             `json:"time" bson:"time,omitempty"`
	BalanceRequester ActorObject        `json:"balance_requester" bson:"balance_requester,omitempty"`
	BalanceOwner     ActorObject        `json:"balance_owner" bson:"balance_owner,omitempty"`
	Access           string             `json:"access" bson:"access"`
	Status           string             `json:"status" bson:"status,omitempty"`
}

// Base interface

func (self RequestAccessBalance) GetDocumentID() primitive.ObjectID {
	return self.ID
}

func (self RequestAccessBalance) SetDocumentID(id primitive.ObjectID) {
	self.ID = id
}

func (self RequestAccessBalance) CollectionName() string {
	return RAB_COLLECTION_NAME
}
