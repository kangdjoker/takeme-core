package usecase

import (
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/usecase/security"
	"github.com/takeme-id/core/utils"
)

func ValidateFormatPIN(pin string) error {
	// validate character PIN
	PINLength := len(pin)
	if PINLength != 6 {
		return utils.ErrorBadRequest(utils.InvalidFormatPIN, "Invalid format PIN")
	}

	return nil
}

func ValidateActorPIN(actor domain.ActorAble, pinEncrypted string) error {

	if actor.IsFaceAsPIN() == false {
		if pinEncrypted == "" {
			return utils.ErrorForbidden()
		}

		pin, err := utils.RSADecrypt(pinEncrypted)
		if err != nil {
			pin, err = utils.RSADecrypDashboard(pinEncrypted)
			if err != nil {
				return utils.ErrorInternalServer(utils.DecryptError, "Decrypt error")
			}
		}

		if pin != actor.GetPIN() {

			a, ok := actor.(domain.User)
			if ok {
				go security.InvalidUserAuth(a)
			} else {
				a, _ := actor.(domain.Corporate)
				go security.InvalidCorporateAuth(a)
			}

			return utils.ErrorForbidden()
		}

		return nil
	} else {
		pin, err := utils.RSADecrypt(pinEncrypted)

		if err != nil {
			return utils.ErrorInternalServer(utils.DecryptError, "Decrypt error")
		}

		if pin != actor.GetTemporaryPIN() {

			a, ok := actor.(domain.User)
			if ok {
				go security.InvalidUserAuth(a)
			} else {
				a, _ := actor.(domain.Corporate)
				go security.InvalidCorporateAuth(a)
			}

			return utils.ErrorForbidden()
		}

		return nil
	}
}

func ValidateAccessBalance(actor domain.ActorAble, balanceID string) error {

	listBalance := actor.GetBalances()

	isFound := false
	for _, element := range listBalance {
		if element.BalanceID.Hex() == balanceID && element.Access == domain.ACCESS_BALANCE_SHARED {
			isFound = true
		}

		if element.BalanceID.Hex() == balanceID && element.Access == domain.ACCESS_BALANCE_OWNER {
			isFound = true
		}
	}

	if isFound == false {
		a, ok := actor.(domain.User)
		if ok {
			go security.InvalidUserAuth(a)
		} else {
			a, _ := actor.(domain.Corporate)
			go security.InvalidCorporateAuth(a)
		}

		return utils.ErrorBadRequest(utils.InvalidBalanceAccess, "Invalid balance access")
	}

	return nil
}

func IsBalanceOwner(actor domain.ActorAble, balanceID string) bool {

	listBalance := actor.GetBalances()

	isFound := false
	for _, element := range listBalance {
		if element.BalanceID.Hex() == balanceID && element.Access == domain.ACCESS_BALANCE_OWNER {
			isFound = true
		}
	}

	return isFound
}

func IsAccessBalanceAlreadyHave(actor domain.ActorAble, balanceID string) bool {

	listBalance := actor.GetBalances()

	isFound := false
	for _, element := range listBalance {
		if element.BalanceID.Hex() == balanceID && element.Access == domain.ACCESS_BALANCE_SHARED {
			isFound = true
		}

		if element.BalanceID.Hex() == balanceID && element.Access == domain.ACCESS_BALANCE_OWNER {
			isFound = true
		}
	}

	return isFound
}

func ValidateIsVerify(actor domain.ActorAble) error {
	if actor.IsVerify() == false {
		return utils.ErrorBadRequest(utils.UpgradeAccountFirst, "Unverified user attempt to transfer")
	}

	return nil
}
