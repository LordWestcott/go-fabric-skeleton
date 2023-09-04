package main

import (
	"appointme/data"
	"appointme/responses"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
	"github.com/stripe/stripe-go/v74"
	up "github.com/upper/db/v4"
)

type StripeCreateSubscriptionRequest struct {
	PriceID                  int    `json:"price_id"`
	StripePaymentMethodID    string `json:"stripe_payment_method_id"`
	PaymentMethodAttached    bool   `json:"payment_method_attached"`
	MakeDefaultPaymentMethod bool   `json:"make_default_payment_method"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	app, err := gofabric.InitApp()
	if err != nil {
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	models := data.New(app.DB)

	// Delete this is you don't need to authenticate the user, and you don't need the user object.
	token := data.Token{}
	user, err := token.AuthenticateTokenFromRequest(r, models.Upper)
	if err != nil {
		if err == up.ErrNoMoreRows || err == up.ErrNilRecord {
			res := responses.NewResponseBase(false, "Invalid token", nil)
			res.WriteJson(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := StripeCreateSubscriptionRequest{}
	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if app.Stripe == nil {
		fmt.Printf("Stripe not configured\n")
		res := responses.NewResponseBase(false, "Stripe not configured", nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Make sure plan exists in db
	plan, err := models.Prices.Get(request.PriceID, models.Upper)
	if err != nil {
		fmt.Printf("Error getting plan from db: %v\n", err)
		res := responses.NewResponseBase(false, "Subscription plan not found.", nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Get the customer object from stripe
	customer, err := app.Stripe.GetCustomer(user.StripeID)
	if err != nil {
		fmt.Printf("Error getting customer from stripe: %v\n", err)
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//If payment method isn't already attached already attach it.
	if !request.PaymentMethodAttached {
		pm, err := app.Stripe.AttachPaymentMethod(user.StripeID, request.StripePaymentMethodID)
		if err != nil {
			fmt.Printf("Error attaching payment method: %v\n", err)
			res := responses.NewResponseBase(false, err.Error(), nil)
			res.WriteJson(w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		request.StripePaymentMethodID = pm.ID
	}

	//Make default payment method if request demands it or if there is none.
	if customer.InvoiceSettings.DefaultPaymentMethod == nil || request.MakeDefaultPaymentMethod {
		updcus, err := app.Stripe.UpdateCustomer(user.StripeID, &stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(request.StripePaymentMethodID),
			},
		})
		if err != nil {
			fmt.Printf("Error updating customer: %v\n", err)
			res := responses.NewResponseBase(false, err.Error(), nil)
			res.WriteJson(w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		customer = updcus
	}

	subparams := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(plan.StripeID),
			},
		},
		Params: stripe.Params{
			Metadata: map[string]string{
				"uID":       strconv.FormatInt(user.ID, 10),
				"ProductID": strconv.FormatInt(plan.ProductID, 10),
				"PriceID":   strconv.FormatInt(plan.ID, 10),
			},
		},
		// TrialPeriodDays: stripe.Int64(14),
	}
	subparams.AddExpand("latest_invoice.payment_intent")

	//Create the subscription
	//If you haven't already, you should create a new product on stripe with recurring billing, add a plan to that.
	sub, err := app.Stripe.CreateSubscription(
		user.StripeID,
		subparams,
	)
	if err != nil {
		fmt.Printf("Error creating subscription: %v\n", err)
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//This is added to the database upon the "customer.subscription.created" webhook event.

	invoice := sub.LatestInvoice

	paymentIntent := invoice.PaymentIntent

	if paymentIntent != nil {
		if paymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
			//TODO handle
		}
	}

	res := responses.NewResponseBase(true, "Subscription Created", sub)

	fmt.Printf("Subscription Data: %v\n", sub)
	w.WriteHeader(http.StatusOK)
	res.WriteJson(w)
}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
