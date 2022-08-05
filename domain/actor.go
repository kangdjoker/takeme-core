package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	ACTOR_TYPE_CORPORATE = CORPORATE_COLLECTION
	ACTOR_TYPE_USER      = USER_COLLECTION
)

type ActorObject struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type      string             `json:"type" bson:"type,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	BalanceID primitive.ObjectID `json:"_" bson:"_,omitempty"`
}

func (self *ActorObject) GetActorID() primitive.ObjectID {
	return self.ID
}

func (self *ActorObject) GetActorType() string {
	return self.Type
}

func (self *ActorObject) GetActorName() string {
	return self.Name
}

func (self *ActorObject) SetActorBalance(balanceID primitive.ObjectID) {
	self.BalanceID = balanceID
}

func (self *ActorObject) GetActorBalance() primitive.ObjectID {
	return self.BalanceID
}

func (self *ActorObject) ToActorObject() ActorObject {
	return ActorObject{
		ID:   self.GetActorID(),
		Type: self.GetActorType(),
		Name: self.GetActorName(),
	}
}

type ActorAble interface {
	GetActorID() primitive.ObjectID
	GetActorType() string
	GetActorName() string
	SetActorBalance(balanceID primitive.ObjectID)
	GetActorBalance() primitive.ObjectID
	GetBalances() []AccessBalance
	GetPIN() string
	GetTemporaryPIN() string
	IsFaceAsPIN() bool
	IsVerify() bool
	ToActorObject() ActorObject
	ToTransactionObject() TransactionObject
}
