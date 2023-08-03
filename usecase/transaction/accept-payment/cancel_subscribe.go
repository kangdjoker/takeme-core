package acceptpayment

import (
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/gateway"
)

func CancelSubscribe(paramLog *basic.ParamLog, subscribeCode string) error {
	gateway := gateway.StripeGateway{}

	err := gateway.CancelSubscribe(paramLog, subscribeCode)
	if err != nil {
		return err
	}

	return nil
}
