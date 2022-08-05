package usecase

import (
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/service"
	"github.com/takeme-id/core/utils"
	"github.com/takeme-id/core/utils/database"
	"go.mongodb.org/mongo-driver/bson"
)

func ActorByID(actorID string) (domain.ActorAble, error) {
	var result domain.ActorAble
	var err error

	result, err = service.UserByIDNoSession(actorID)
	if err != nil {
		return result, err
	}

	if result.GetActorType() == "" {
		result, err = service.CorporateByIDNoSession(actorID)
		if err != nil {
			return result, err
		}
	}

	if result.GetActorType() == "" {
		return result, utils.ErrorBadRequest(utils.AccountNotFound, "Account not found")
	}

	return result, nil
}

func ActorAddBalance(actor domain.ActorAble, newAccessBalance domain.AccessBalance) error {
	accessBalance := actor.GetBalances()
	accessBalance = append(accessBalance, newAccessBalance)

	query := bson.M{"$set": bson.M{"list_balance": accessBalance}}
	err := database.UpdateQuery(actor.GetActorType(), actor.GetActorID(), query)
	if err != nil {
		return err
	}

	return nil
}

func ActorRemoveBalance(actor domain.ActorAble, balanceID string) error {
	accessBalance := actor.GetBalances()

	var newAccessBalance []domain.AccessBalance
	for _, element := range accessBalance {
		if balanceID != element.BalanceID.Hex() {
			newAccessBalance = append(newAccessBalance, element)
		}
	}

	accessBalance = newAccessBalance

	query := bson.M{"$set": bson.M{"list_balance": accessBalance}}
	err := database.UpdateQuery(actor.GetActorType(), actor.GetActorID(), query)
	if err != nil {
		return err
	}

	return nil
}

func ActorObjectToActor(actor domain.ActorObject) (domain.ActorAble, error) {

	collection := actor.Type
	ID := actor.GetActorID()
	var result domain.ActorAble
	var err error

	if collection == domain.USER_COLLECTION {
		result, err = service.UserByIDNoSession(ID.Hex())
		if err != nil {
			return result, err
		}
	} else {
		result, err = service.CorporateByIDNoSession(ID.Hex())
		if err != nil {
			return result, err
		}
	}

	return result, nil
}
