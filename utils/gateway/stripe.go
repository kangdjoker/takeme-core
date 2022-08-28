package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils"
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

func (gateway StripeGateway) CallbackAcceptPaymentCard(w http.ResponseWriter, r *http.Request) (string, int, domain.Card, string, error) {

	log.Info("------------------------ Stripe hit callback card ------------------------")

	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload := r.Context().Value("payload").([]byte)

	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := os.Getenv("STRIPE_SECRET")
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		return "", 0, domain.Card{}, "", utils.ErrorBadRequest(utils.InvalidSignature, "Invalid signature stripe callback")
	}

	if event.Type != "payment_intent.succeeded" {
		return "", 0, domain.Card{}, "", utils.ErrorBadRequest(utils.InvalidRequestPayload, "Invalid payload stripe callback")
	}

	var paymentIntent stripe.PaymentIntent
	err = json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return "", 0, domain.Card{}, "", utils.ErrorBadRequest(utils.InvalidRequestPayload, "Invalid payload stripe callback")
	}

	card := domain.Card{
		AccountNumber: "**** **** **** " + paymentIntent.PaymentMethod.Card.Last4,
		ExpMonth:      strconv.Itoa(int(paymentIntent.PaymentMethod.Card.ExpMonth)),
		ExpYear:       strconv.Itoa(int(paymentIntent.PaymentMethod.Card.ExpYear)),
		Network:       string(paymentIntent.PaymentMethod.Card.Brand),
	}

	amount := int(paymentIntent.Amount)
	reference := paymentIntent.PaymentMethod.ID
	balanceID := paymentIntent.Description

	return balanceID, amount, card, reference, nil
}
