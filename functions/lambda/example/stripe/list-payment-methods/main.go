package main

import (
	"appointme/data"
	"appointme/responses"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
	"github.com/stripe/stripe-go/v74"
	up "github.com/upper/db/v4"
)

// Saves a card on the stripe customer account with a setup intent
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

	// Create a setup intent

	if user.StripeID == "" {
		res := responses.NewResponseBase(false, "User has no stripe ID", nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusUnauthorized)
	}

	pmIter := app.Stripe.ListPaymentMethods(&stripe.PaymentMethodListParams{
		Customer: stripe.String(user.StripeID),
		Type:     stripe.String(string(stripe.PaymentMethodTypeCard)),
	})

	paymentMethods := []*stripe.PaymentMethod{}
	for pmIter.Next() {
		if pmIter.Err() != nil {
			res := responses.NewResponseBase(false, pmIter.Err().Error(), nil)
			res.WriteJson(w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		paymentMethods = append(paymentMethods, pmIter.PaymentMethod())
	}

	res := responses.NewResponseBase(true, "Payment Methods", paymentMethods)
	w.WriteHeader(http.StatusOK)
	res.WriteJson(w)
}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
