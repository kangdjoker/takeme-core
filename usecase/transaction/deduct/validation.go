package deduct

import (
	"os"
	"strconv"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
)

func validateMinimum(transaction domain.Transaction) error {
	minimum, _ := strconv.Atoi(os.Getenv("MINIMUM_TRANSFER_AMOUNT"))
	if transaction.Amount < minimum {
		return utils.ErrorBadRequest(utils.MinimumAmountTransaction, "Transaction under minimum")
	}

	return nil
}

func validateMaximum(transaction domain.Transaction) error {
	maximum, _ := strconv.Atoi(os.Getenv("MAXIMUM_TRANSFER_AMOUNT"))
	if transaction.Amount > maximum {
		return utils.ErrorBadRequest(utils.MaximumAmountTransaction, "Transaction reach maximum")
	}

	return nil
}

func validateCurrency(from domain.Balance, to domain.Balance) error {
	if from.Currency != to.Currency {
		return utils.ErrorBadRequest(utils.CurrencyError, "Transaction cross currency")
	}

	return nil
}
