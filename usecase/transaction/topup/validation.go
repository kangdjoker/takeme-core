package topup

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

func validateCurrency(paramLog *basic.ParamLog, from domain.Balance, to domain.Balance) error {
	if from.Currency != to.Currency {
		return utils.ErrorBadRequest(paramLog, utils.CurrencyError, "Transaction cross currency")
	}

	return nil
}
