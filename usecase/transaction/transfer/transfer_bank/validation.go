package transfer_bank

import (
	"os"
	"strconv"

	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils"
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

func validateCurrency(transaction domain.Transaction, corporate domain.Corporate) error {
	if transaction.Currency != "idr" {
		return utils.ErrorBadRequest(utils.OnlySupportOnIDR, "Transaction cross currency")
	}

	return nil
}
