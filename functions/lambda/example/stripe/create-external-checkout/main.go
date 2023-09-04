package main

import (
	"appointme/responses"
	"encoding/json"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
	"github.com/stripe/stripe-go/v74"
)

type StripeExternalCheckoutResponse struct {
	SessionID string `json:"session_id"`
}

type StripeExternalCheckoutRequest struct {
	Currency string `json:"currency"`
	Product  struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
		Price       int64    `json:"price"`
	} `json:"product"`
	Quantity           int64    `json:"quantity"`
	SuccessURL         string   `json:"success_url"`
	CancelURL          string   `json:"cancel_url"`
	PaymentMethodTypes []string `json:"payment_method_types"`
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
	// models := data.New(app.DB)

	request := StripeExternalCheckoutRequest{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := stripe.CheckoutSessionParams{}
	//required
	params.Mode = stripe.String(string(stripe.CheckoutSessionModePayment)) // or stripe.CheckoutSessionModeSetup or stripe.CheckoutSessionModeSubscription
	params.LineItems = []*stripe.CheckoutSessionLineItemParams{
		{
			// Price:   stripe.String("price_000000000000000000"), //or reference price here.
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(string(request.Currency)),
				//Currency: stripe.String(string(stripe.CurrencyUSD)), //or set on server side
				//Product: stripe.String("prod_00000000000000"), //or reference product here.
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(request.Product.Name),
					Description: stripe.String(request.Product.Description),
					Images:      stripe.StringSlice(request.Product.Images),
				},
				// ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				// 	Name:        stripe.String("T-shirt"),
				// 	Description: stripe.String("Comfortable cotton t-shirt"),
				// 	Images: []*string{
				// 		stripe.String("https://plus.unsplash.com/premium_photo-1673125287084-e90996bad505?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=987&q=80"),
				// 	},
				// }, // or set on server side.
				UnitAmount: stripe.Int64(request.Product.Price), // $20.00
				// UnitAmount: stripe.Int64(2000), // $20.00 - or set on server side.
				// Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{ // for subscriptions
				// 	Interval: stripe.String(string("month")),
				// },
			},
			Quantity: stripe.Int64(request.Quantity),
		},
	}

	params.SuccessURL = stripe.String(request.SuccessURL)
	// params.SuccessURL = stripe.String(frontend + "/examples/stripe/redirects/payment-success?session_id={CHECKOUT_SESSION_ID}") // or set on server side
	params.PaymentMethodTypes = stripe.StringSlice(request.PaymentMethodTypes)
	// params.PaymentMethodTypes = stripe.StringSlice([]string{
	// 	"card",
	// }) //or set on server side

	//not required
	params.CancelURL = stripe.String(request.SuccessURL)
	// params.CancelURL = stripe.String(frontend + "/examples/stripe/redirects/payment-cancel?session_id={CHECKOUT_SESSION_ID}") // or set on server side

	session, error := app.Stripe.CreateCheckoutSession(&params)
	if error != nil {
		res := responses.NewResponseBase(false, error.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := StripeExternalCheckoutResponse{
		SessionID: session.ID,
	}

	res := responses.NewResponseBase(true, "External Checkout Created", response)
	w.WriteHeader(http.StatusOK)
	res.WriteJson(w)

}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
