package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/customer"
	"github.com/stripe/stripe-go/v73/invoice"
	"github.com/stripe/stripe-go/v73/paymentintent"
	"github.com/stripe/stripe-go/v73/paymentmethod"
	"github.com/stripe/stripe-go/v73/price"
	"github.com/stripe/stripe-go/v73/product"
	"github.com/stripe/stripe-go/v73/subscription"
	"github.com/stripe/stripe-go/v73/webhook"
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

func (gateway StripeGateway) ChargeCard(balanceID string, amount int, returnURL string, card domain.Card, externalID string) (string, string, error) {
	stripe.Key = os.Getenv("STRIPE_SECRET")

	expM, err := strconv.ParseInt(card.ExpMonth, 10, 64)
	expY, err := strconv.ParseInt(card.ExpYear, 10, 64)
	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   &card.AccountNumber,
			ExpMonth: &expM,
			ExpYear:  &expY,
			CVC:      &card.CVC,
		},
		Type: stripe.String("card"),
	}
	pm, _ := paymentmethod.New(params)

	paymentMethodID := pm.ID

	reference := balanceID

	params2 := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		UseStripeSDK:  stripe.Bool(false),
		PaymentMethod: &paymentMethodID,
	}

	params2.AddMetadata("reference", reference)
	params2.AddMetadata("external_id", externalID)

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

func (gateway StripeGateway) ChargeCardSubscribe(balanceID string, amount int, returnURL string, card domain.Card, externalID string, interval string) (
	string, string, string, error) {
	stripe.Key = os.Getenv("STRIPE_SECRET")
	reference := balanceID

	expM, err := strconv.ParseInt(card.ExpMonth, 10, 64)
	expY, err := strconv.ParseInt(card.ExpYear, 10, 64)
	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   &card.AccountNumber,
			ExpMonth: &expM,
			ExpYear:  &expY,
			CVC:      &card.CVC,
		},
		Type: stripe.String("card"),
	}
	params.AddMetadata("reference", reference)
	pm, err := paymentmethod.New(params)
	if err != nil {
		return "", "", "", utils.ErrorInternalServer(utils.StripeAPICallFail, "Stripe API call fail")
	}

	paymentMethodID := pm.ID

	params2 := &stripe.ProductParams{
		Name: stripe.String("Gold Special"),
	}
	pro, err := product.New(params2)

	productID := pro.ID

	priceParam := &stripe.PriceParams{
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Product:  &productID,
		Recurring: &stripe.PriceRecurringParams{
			Interval: &interval,
		},
		UnitAmount: stripe.Int64(2000),
	}
	pr, _ := price.New(priceParam)
	priceID := pr.ID

	custParams := &stripe.CustomerParams{
		Description:   stripe.String(balanceID),
		PaymentMethod: &paymentMethodID,
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: &paymentMethodID,
		},
	}

	c, err := customer.New(custParams)
	if err != nil {
		return "", "", "", utils.ErrorInternalServer(utils.StripeAPICallFail, "Stripe API call fail")
	}

	customerID := c.ID

	behaviour := "allow_incomplete"
	subParam := &stripe.SubscriptionParams{
		Customer: &customerID,
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: &priceID,
			},
		},
		PaymentBehavior: &behaviour,
	}
	s, err := subscription.New(subParam)

	subscriptionID := s.ID
	invoiceID := s.LatestInvoice.ID

	in, err := invoice.Get(
		invoiceID,
		nil,
	)
	if err != nil {
		return "", "", "", utils.ErrorInternalServer(utils.StripeAPICallFail, "Stripe API call fail")
	}

	paymentIntentID := in.PaymentIntent.ID

	pi, _ := paymentintent.Get(
		paymentIntentID,
		nil,
	)

	status := CHARGE_CARD_STATUS_PENDING
	authURL := ""
	if string(pi.Status) == "requires_action" {
		authURL = pi.NextAction.RedirectToURL.URL
	} else {
		status = CHARGE_CARD_STATUS_COMPLETED
	}

	return status, authURL, subscriptionID, nil
}

func (gateway StripeGateway) CancelSubscribe(subsID string) error {
	stripe.Key = os.Getenv("STRIPE_SECRET")

	_, err := subscription.Cancel(
		subsID,
		nil,
	)
	if err != nil {
		return utils.ErrorInternalServer(utils.StripeAPICallFail, "Stripe API call fail")
	}

	return nil
}

func (gateway StripeGateway) CallbackAcceptPaymentCard(w http.ResponseWriter, r *http.Request) (string, int, domain.Card, string, string, error) {

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
		return "", 0, domain.Card{}, "", "", utils.ErrorBadRequest(utils.InvalidSignature, "Invalid signature stripe callback")
	}

	if event.Type != "payment_intent.succeeded" {
		return "", 0, domain.Card{}, "", "", nil
	}

	var paymentIntent stripe.PaymentIntent
	err = json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return "", 0, domain.Card{}, "", "", utils.ErrorBadRequest(utils.InvalidRequestPayload, "Invalid payload stripe callback")
	}

	card := domain.Card{
		AccountNumber: "**** **** **** " + paymentIntent.Charges.Data[0].PaymentMethodDetails.Card.Last4,
		ExpMonth:      strconv.Itoa(int(paymentIntent.Charges.Data[0].PaymentMethodDetails.Card.ExpMonth)),
		ExpYear:       strconv.Itoa(int(paymentIntent.Charges.Data[0].PaymentMethodDetails.Card.ExpYear)),
		Network:       string(paymentIntent.Charges.Data[0].PaymentMethodDetails.Card.Brand),
	}

	amount := int(paymentIntent.Amount)
	reference := paymentIntent.ID
	balanceID := paymentIntent.Metadata["reference"]
	externalID := paymentIntent.Metadata["external_id"]

	if balanceID == "" {
		balanceID = paymentIntent.PaymentMethod.Metadata["reference"]
	}

	if externalID == "" {
		externalID = paymentIntent.PaymentMethod.Metadata["external_id"]
	}

	return balanceID, amount, card, reference, externalID, nil
}
