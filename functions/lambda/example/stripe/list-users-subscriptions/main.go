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

	// DO STUFF HERE...

	iter := app.Stripe.ListSubscriptions(&stripe.SubscriptionListParams{
		Customer: stripe.String(user.StripeID),
	})

	subscriptions := []*stripe.Subscription{}
	for iter.Next() {
		if iter.Err() != nil {
			res := responses.NewResponseBase(false, iter.Err().Error(), nil)
			res.WriteJson(w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		subscriptions = append(subscriptions, iter.Subscription())
	}

	res := responses.NewResponseBase(true, "User Subscriptions", subscriptions)
	w.WriteHeader(http.StatusOK)
	res.WriteJson(w)

}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
