package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel interface {
	GetDocumentID() primitive.ObjectID
	SetDocumentID(id primitive.ObjectID)
	CollectionName() string
}

type Audit struct {
	CreatedTime string `json:"created_time" bson:"created_time,omitempty"`
	UpdatedTime string `json:"updated_time" bson:"updated_time,omitempty"`
}
