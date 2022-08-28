package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/webhook"
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils"
)

const (
	CHARGE_CARD_STATUS_COMPLETED = "Completed"
	CHARGE_CARD_STATUS_PENDING   = "Pending"
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

func (gateway StripeGateway) ChargeCard(balanceID string, returnURL string, card domain.Card) (string, string, error) {
	stripe.Key = os.Getenv("STRIPE_SECRET")

	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   &card.AccountNumber,
			ExpMonth: &card.ExpMonth,
			ExpYear:  &card.ExpYear,
			CVC:      &card.CVC,
		},
		Type: stripe.String("card"),
	}
	pm, _ := paymentmethod.New(params)

	paymentMethodID := pm.ID

	reference := balanceID

	params2 := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(2000),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: &reference,
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		UseStripeSDK:  stripe.Bool(false),
		PaymentMethod: &paymentMethodID,
	}

	pi, err := paymentintent.New(params2)

	if err != nil {
		return "", "", utils.ErrorInternalServer(utils.StripeAPICallFail, "Stripe API call fail")
	}

	params3 := &stripe.PaymentIntentConfirmParams{
		UseStripeSDK: stripe.Bool(false),
		ReturnURL:    &returnURL,
	}

	pi2, err := paymentintent.Confirm(
		pi.ID,
		params3,
	)

	if err != nil {
		return "", "", utils.ErrorInternalServer(utils.StripeAPICallFail, "Stripe API call fail")
	}

	status := CHARGE_CARD_STATUS_PENDING
	authURL := ""
	if string(pi2.Status) == "requires_action" {
		authURL = pi2.NextAction.RedirectToURL.URL
	} else {
		status = CHARGE_CARD_STATUS_COMPLETED
	}

	return status, authURL, nil
}

func (gateway StripeGateway) CallbackAcceptPaymentCard(w http.ResponseWriter, r *http.Request) (string, int, domain.Card, string, error) {

	log.Info("------------------------ Stripe hit callback card ------------------------")

	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload := r.Context().Value("payload").([]byte)

	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		return "", 0, domain.Card{}, "", utils.ErrorBadRequest(utils.InvalidSignature, "Invalid signature stripe callback")
	}

	if event.Type != "payment_intent.succeeded" {
		return "", 0, domain.Card{}, "", nil
	}

	var paymentIntent stripe.PaymentIntent
	err = json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return "", 0, domain.Card{}, "", utils.ErrorBadRequest(utils.InvalidRequestPayload, "Invalid payload stripe callback")
	}

	card := domain.Card{
		AccountNumber: "**** **** **** " + paymentIntent.Charges.Data[0].PaymentMethodDetails.Card.Last4,
		ExpMonth:      strconv.Itoa(int(paymentIntent.Charges.Data[0].PaymentMethodDetails.Card.ExpMonth)),
		ExpYear:       strconv.Itoa(int(paymentIntent.Charges.Data[0].PaymentMethodDetails.Card.ExpYear)),
		Network:       string(paymentIntent.Charges.Data[0].PaymentMethodDetails.Card.Brand),
	}

	amount := int(paymentIntent.Amount)
	reference := paymentIntent.PaymentMethod.ID
	balanceID := paymentIntent.Description

	return balanceID, amount, card, reference, nil
}
