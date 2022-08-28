package gateway

import (
	"net/http"

	"github.com/takeme-id/core/domain"
)

const (
	Xendit  = "A"
	Fusindo = "B"
	Sprint  = "C"
	OY      = "D"
	MMBC    = "E"
	Stripe  = "F"
)

type Gateway interface {
	Name() string
	CreateVA(balanceID string, nameVA string, bankCode string) (string, error)
	CallbackVA(w http.ResponseWriter, r *http.Request) (string, int, domain.Bank, string, error)
	CreateTransfer(transaction domain.Transaction) (string, error)
	CallbackTransfer(w http.ResponseWriter, r *http.Request) (string, string, string, error)
	Inquiry(bankCode string, accountNumber string) (string, error)
}
