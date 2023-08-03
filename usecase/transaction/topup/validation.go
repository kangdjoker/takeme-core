package topup

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

func validateCurrency(paramLog *basic.ParamLog, incomeCurrency string, balance domain.Balance) error {
	if incomeCurrency != balance.Currency {
		return utils.ErrorBadRequest(paramLog, utils.CurrencyError, "Transaction cross currency")
	}

	return nil
}
