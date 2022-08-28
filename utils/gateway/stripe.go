package gateway

import (
	"net/http"

	"github.com/takeme-id/core/domain"
)

type StripeGateway struct {
}

func (gateway StripeGateway) Name() string {
	return Stripe
}

func (gateway StripeGateway) CreateVA(balanceID string, nameVA string, bankCode string) (string, error) {
	return "", nil
}
func (gateway StripeGateway) CallbackVA(w http.ResponseWriter, r *http.Request) (string, int, domain.Bank, string, error) {
	return "", 0, domain.Bank{}, "", nil
}

func (gateway StripeGateway) CreateTransfer(transaction domain.Transaction) (string, error) {
	return "", nil
}
func (gateway StripeGateway) CallbackTransfer(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
	return "", "", "", nil
}
func (gateway StripeGateway) Inquiry(bankCode string, accountNumber string) (string, error) {
	return "", nil
}
