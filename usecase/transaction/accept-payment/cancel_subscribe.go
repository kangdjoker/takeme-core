package acceptpayment

import (
	"github.com/takeme-id/core/utils/gateway"
)

func CancelSubscribe(subscribeCode string) error {
	gateway := gateway.StripeGateway{}

	err := gateway.CancelSubscribe(subscribeCode)
	if err != nil {
		return err
	}

	return nil
}
