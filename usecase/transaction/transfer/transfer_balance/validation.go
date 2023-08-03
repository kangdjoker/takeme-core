package transfer_balance

import (
	"os"
	"strconv"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

func validateMinimum(paramLog *basic.ParamLog, transaction domain.Transaction) error {
	minimum, _ := strconv.Atoi(os.Getenv("MINIMUM_TRANSFER_AMOUNT"))
	if transaction.Amount < minimum {
		return utils.ErrorBadRequest(paramLog, utils.MinimumAmountTransaction, "Transaction under minimum")
	}

	return nil
}

func validateMaximum(paramLog *basic.ParamLog, transaction domain.Transaction) error {
	maximum, _ := strconv.Atoi(os.Getenv("MAXIMUM_TRANSFER_AMOUNT"))
	if transaction.Amount > maximum {
		return utils.ErrorBadRequest(paramLog, utils.MaximumAmountTransaction, "Transaction reach maximum")
	}

	return nil
}

func validateCurrency(paramLog *basic.ParamLog, from domain.Balance, to domain.Balance) error {
	if from.Currency != to.Currency {
		return utils.ErrorBadRequest(paramLog, utils.CurrencyError, "Transaction cross currency")
	}

	return nil
}
