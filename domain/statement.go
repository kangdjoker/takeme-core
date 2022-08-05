package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const STATEMENT_COLLECTION_NAME string = "statement"
const (
	STATEMENT_TYPE_FEE         = "FEE"
	STATEMENT_TYPE_TRANSACTION = "TRANSACTION"
)

type Statement struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BalanceID   primitive.ObjectID `json:"balance_id" bson:"balance_id,omitempty"`
	Time        string             `json:"time" bson:"time,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Reference   string             `json:"reference" bson:"reference,omitempty"`
	Withdraw    int                `json:"withdraw" bson:"withdraw"`
	Deposit     int                `json:"deposit" bson:"deposit"`
	Balance     int                `json:"balance" bson:"balance"`
	Type        string             `json:"type" bson:"type,omitempty"`
}

// Base interface

func (self Statement) GetDocumentID() primitive.ObjectID {
	return self.ID
}

func (self Statement) SetDocumentID(id primitive.ObjectID) {
	self.ID = id
}

func (self Statement) CollectionName() string {
	return STATEMENT_COLLECTION_NAME
}
