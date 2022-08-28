package acceptpayment

import (
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils"
)

func validateCurrency(incomeCurrency string, balance domain.Balance) error {
	if incomeCurrency != balance.Currency {
		return utils.ErrorBadRequest(utils.CurrencyError, "Transaction cross currency")
	}

	return nil
}
