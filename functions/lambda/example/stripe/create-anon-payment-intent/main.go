package main

import (
	"appointme/responses"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
	"github.com/stripe/stripe-go/v74"
)

type StripeAnonPaymentIntentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type StripeAnonPaymentIntentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	ClientSecret string `json:"client_secret"`
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

	request := StripeAnonPaymentIntentRequest{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Request: %+v\n", request)

	paymentintent, err := app.Stripe.CreatePaymentIntent(&stripe.PaymentIntentParams{
		Amount:   stripe.Int64(request.Amount),
		Currency: stripe.String(request.Currency),
	})
	if err != nil {
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := StripeAnonPaymentIntentResponse{
		ID:           paymentintent.ID,
		Status:       string(paymentintent.Status),
		ClientSecret: paymentintent.ClientSecret,
	}

	res := responses.NewResponseBase(true, "Payment Intent Created", response)
	w.WriteHeader(http.StatusOK)
	res.WriteJson(w)

}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
