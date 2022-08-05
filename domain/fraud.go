package domain

import (
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mitigation
const (
	USER_FAILED_ATTEMPT      = "User failed authenticate"
	CORPORATE_FAILED_ATTEMPT = "Corporate failed authenticate"
	USER_LOCKED              = "User locked"
	CORPORATE_LOCKED         = "Corporate locked"
	TRANSACTION_CANCELED     = "Transaction canceled because detected as identycal transaction"
)

const FRAUD_COLLECTION string = "fraud"

type Fraud struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Actor       ActorObject        `json:"actor" bson:"actor,omitempty"`
	Time        string             `json:"time" bson:"time,omitempty"`
}

func CreateFraud(description string, actor ActorAble, actorType string) Fraud {
	return Fraud{
		Description: description,
		Time:        time.Now().Format(os.Getenv("TIME_FORMAT")),
		Actor:       actor.ToActorObject(),
	}
}

// Interface for mongo document result
func (domain *Fraud) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *Fraud) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (domain *Fraud) CollectionName() string {
	return FRAUD_COLLECTION
}
