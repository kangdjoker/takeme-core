package usecase

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/usecase/security"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

func ValidateFormatPIN(paramLog *basic.ParamLog, pin string) error {
	// validate character PIN
	PINLength := len(pin)
	if PINLength != 6 {
		return utils.ErrorBadRequest(&basic.ParamLog{}, utils.InvalidFormatPIN, "Invalid format PIN")
	}

	return nil
}

func ValidateActorPIN(paramLog *basic.ParamLog, actor domain.ActorAble, pinEncrypted string) error {

	if actor.IsFaceAsPIN() == false {
		if pinEncrypted == "" {
			basic.LogInformation(paramLog, "error.pinEncrypted")
			return utils.ErrorForbidden(paramLog)
		}

		pin, err := utils.RSADecrypt(pinEncrypted)
		if err != nil {
			pin, err = utils.RSADecrypDashboard(paramLog, pinEncrypted)
			if err != nil {
				return utils.ErrorInternalServer(paramLog, utils.DecryptError, "Decrypt error")
			}
		}

		if pin != actor.GetPIN() {

			a, ok := actor.(domain.User)
			if ok {
				go security.InvalidUserAuth(paramLog, a)
			} else {
				a, _ := actor.(domain.Corporate)
				go security.InvalidCorporateAuth(paramLog, a)
			}

			basic.LogInformation(paramLog, "error.pin != actor.GetPIN()")
			return utils.ErrorForbidden(paramLog)
		}

		return nil
	} else {
		pin, err := utils.RSADecrypt(pinEncrypted)

		if err != nil {
			basic.LogInformation(paramLog, "error.RSADecrypt")
			return utils.ErrorInternalServer(paramLog, utils.DecryptError, "Decrypt error")
		}

		if pin != actor.GetTemporaryPIN() {

			a, ok := actor.(domain.User)
			if ok {
				go security.InvalidUserAuth(paramLog, a)
			} else {
				a, _ := actor.(domain.Corporate)
				go security.InvalidCorporateAuth(paramLog, a)
			}
			basic.LogInformation(paramLog, "error.pin != actor.GetTemporaryPIN()")
			return utils.ErrorForbidden(paramLog)
		}

		return nil
	}
}

func ValidateAccessBalance(paramLog *basic.ParamLog, actor domain.ActorAble, balanceID string) error {

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
			go security.InvalidUserAuth(paramLog, a)
		} else {
			a, _ := actor.(domain.Corporate)
			go security.InvalidCorporateAuth(paramLog, a)
		}

		return utils.ErrorBadRequest(paramLog, utils.InvalidBalanceAccess, "Invalid balance access")
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

func ValidateIsVerify(paramLog *basic.ParamLog, actor domain.ActorAble) error {
	if actor.IsVerify() == false {
		return utils.ErrorBadRequest(paramLog, utils.UpgradeAccountFirst, "Unverified user attempt to transfer")
	}

	return nil
}
