package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const LOG_COLLECTION string = "log"

const (
	ACCESS_LOG_VIEW_ONLY = "View Only"
	ACCESS_LOG_SHARED    = "Shared"
	ACCESS_LOG_OWNER     = "Owner"
)

type Log struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TimeCreate string             `json:"time_create" bson:"time_create,omitempty"`
	Tag        string             `json:"tag" bson:"tag,omitempty"`
	Data       interface{}        `json:"data" bson:"data"`
}

func (domain *Log) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *Log) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (domain *Log) CollectionName() string {
	return LOG_COLLECTION
}
