package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const CALLBACK_HISTORY_COLLECTION string = "callback_history"

type CallbackHistory struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Time            string             `json:"time" bson:"time,omitempty"`
	URL             string             `json:"url" bson:"url,omitempty"`
	TransactionCode string             `json:"transaction_code" bson:"transaction_code,omitempty"`
	RequestBody     string             `json:"request_body" bson:"request_body,omitempty"`
	ResponseBody    string             `json:"response_body" bson:"response_body,omitempty"`
	ResponseStatus  string             `json:"response_status" bson:"response_status,omitempty"`
}

// Interface for mongo document result
func (domain *CallbackHistory) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *CallbackHistory) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (domain *CallbackHistory) CollectionName() string {
	return CALLBACK_HISTORY_COLLECTION
}
