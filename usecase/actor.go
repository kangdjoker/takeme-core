package usecase

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/bson"
)

func ActorByID(paramLog *basic.ParamLog, actorID string) (domain.ActorAble, error) {
	var result domain.ActorAble
	var err error

	result, err = service.UserByIDNoSession(paramLog, actorID)
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
		return result, utils.ErrorBadRequest(paramLog, utils.AccountNotFound, "Account not found")
	}

	return result, nil
}

func ActorAddBalance(paramLog *basic.ParamLog, actor domain.ActorAble, newAccessBalance domain.AccessBalance) error {
	accessBalance := actor.GetBalances()
	accessBalance = append(accessBalance, newAccessBalance)

	query := bson.M{"$set": bson.M{"list_balance": accessBalance}}
	err := database.UpdateQuery(paramLog, actor.GetActorType(), actor.GetActorID(), query)
	if err != nil {
		return err
	}

	return nil
}

func ActorRemoveBalance(paramLog *basic.ParamLog, actor domain.ActorAble, balanceID string) error {
	accessBalance := actor.GetBalances()

	var newAccessBalance []domain.AccessBalance
	for _, element := range accessBalance {
		if balanceID == element.BalanceID.Hex() && element.Access == domain.ACCESS_BALANCE_OWNER {
			return utils.ErrorBadRequest(paramLog, utils.InvalidLevelAccessRevoke, "Invalid level access revoke")
		}

		if balanceID != element.BalanceID.Hex() && element.Access != domain.ACCESS_BALANCE_OWNER {
			newAccessBalance = append(newAccessBalance, element)
		}
	}

	accessBalance = newAccessBalance

	query := bson.M{"$set": bson.M{"list_balance": accessBalance}}
	err := database.UpdateQuery(paramLog, actor.GetActorType(), actor.GetActorID(), query)
	if err != nil {
		return err
	}

	return nil
}

func ActorObjectToActor(paramLog *basic.ParamLog, actor domain.ActorObject) (domain.ActorAble, error) {

	collection := actor.Type
	ID := actor.GetActorID()
	var result domain.ActorAble
	var err error

	if collection == domain.USER_COLLECTION {
		result, err = service.UserByIDNoSession(paramLog, ID.Hex())
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
