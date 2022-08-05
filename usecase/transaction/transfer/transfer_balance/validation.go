package transfer_balance

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
