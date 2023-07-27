package acceptpayment

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
)

func validateCurrency(incomeCurrency string, balance domain.Balance) error {
	if incomeCurrency != balance.Currency {
		return utils.ErrorBadRequest(utils.CurrencyError, "Transaction cross currency")
	}

	return nil
}
