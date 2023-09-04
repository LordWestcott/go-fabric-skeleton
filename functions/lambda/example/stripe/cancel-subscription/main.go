package main

import (
	"appointme/data"
	"appointme/responses"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
	"github.com/stripe/stripe-go/v74"
	up "github.com/upper/db/v4"
)

type StripeCreateSubscriptionRequest struct {
	StripeSubscriptionID string `json:"subscription_id"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	app, err := gofabric.InitApp()
	if err != nil {
		fmt.Printf("Error initializing app: %v\n", err)
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
		fmt.Printf("Error authenticating token: %v\n", err)
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

	defer r.Body.Close()
	request := StripeCreateSubscriptionRequest{}
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sub, err := models.Subscriptions.GetByStripeSubscriptionID(request.StripeSubscriptionID, models.Upper)
	if err != nil {
		fmt.Printf("Error getting subscription: %v\n", err)
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sub.UID != user.ID {
		fmt.Printf("Subscription doesn't belong to user\n")
		res := responses.NewResponseBase(false, "Invalid Subscription", nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//This shall cancel the subscription at the end of the period
	//Once this happens the webhook will be called - subscription.deleted.
	_, err = app.Stripe.UpdateSubscription(sub.StripeSubscriptionID, &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	})
	if err != nil {
		fmt.Printf("Error cancelling subscription on stripe: %v\n", err)
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := responses.NewResponseBase(true, "Subscription Cancelled", nil)
	w.WriteHeader(http.StatusOK)
	res.WriteJson(w)

}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
