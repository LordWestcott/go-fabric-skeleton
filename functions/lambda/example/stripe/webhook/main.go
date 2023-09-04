package main

import (
	"appointme/data"
	"appointme/responses"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
)

// Use this to fulfill stuff upon a webhook event or handle a webhook event

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

	secret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if secret == "" {
		fmt.Println("STRIPE_WEBHOOK_SECRET not set")
		return
	}

	//Verify with stripe signing secret that the request is from stripe.
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	event, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"),
		secret, webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//Handle the event by Type.
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("PaymentIntent was successful!")
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("PaymentMethod was attached to a Customer!")

	case "invoice.payment_failed":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// //Get user if you need them.
		// sCustomer := invoice.Customer
		// uID, ok := sCustomer.Metadata["uID"]
		// if !ok {
		// 	fmt.Fprintf(os.Stderr, "No user with uID %s\n", uID)
		// 	w.WriteHeader(http.StatusNotFound)
		// 	return
		// }
		// userID, err := strconv.Atoi(uID)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Error parsing uID %s\n", uID)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }
		// user, err := models.Users.Get(userID, models.Upper)

		//Make subscription status "past_due"
		//If subscription is "past_due" you should update the frontend to prompt the user to update their payment method.
		stripeSub := invoice.Subscription
		if stripeSub != nil {
			sub, err := models.Subscriptions.GetByStripeSubscriptionID(stripeSub.ID, models.Upper)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting subscription with stripe ID %s\n", stripeSub.ID)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			sub.Status = string(stripe.SubscriptionStatusPastDue)
			if err := models.Subscriptions.Update(*sub, models.Upper); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating subscription with stripe ID %s\n", stripeSub.ID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			//Handle products, if custom stuff is needed here.
			// products := invoice.Lines.Data
			// for _, product := range products {
			// }
		}
		w.WriteHeader(http.StatusOK)

	case "invoice.payment_succeeded":
		//If subscription exists in relation to this invoice, update the status to active.
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		stripeSub := invoice.Subscription
		if stripeSub != nil {
			sub, err := models.Subscriptions.GetByStripeSubscriptionID(stripeSub.ID, models.Upper)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting subscription with stripe ID %s\n", stripeSub.ID)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			sub.Status = string(stripe.SubscriptionStatusActive)
			if err := models.Subscriptions.Update(*sub, models.Upper); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating subscription with stripe ID %s\n", stripeSub.ID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}
	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		uID, ok := subscription.Metadata["uID"]
		if !ok {
			fmt.Fprintf(os.Stderr, "No user with uID %s\n", uID)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		userID, err := strconv.Atoi(uID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing uID %s\n", uID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pID, ok := subscription.Metadata["ProductID"]
		if !ok {
			fmt.Fprintf(os.Stderr, "No product with ID %s\n", pID)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		productID, err := strconv.Atoi(pID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing product ID %s\n", pID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		prID, ok := subscription.Metadata["PriceID"]
		if !ok {
			fmt.Fprintf(os.Stderr, "No price with ID %s\n", prID)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		priceID, err := strconv.Atoi(prID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing price ID %s\n", prID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sub := data.Subscription{
			UID:                  int64(userID),
			ProductID:            int64(productID),
			PriceID:              int64(priceID),
			StripeSubscriptionID: subscription.ID,
			StripeCustomerID:     subscription.Customer.ID,
			StripePriceID:        subscription.Items.Data[0].Price.ID,
			Status:               string(subscription.Status),
		}

		_, err = models.Subscriptions.Insert(&sub, models.Upper)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating subscription with stripe ID %s\n", subscription.ID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sub, err := models.Subscriptions.GetByStripeSubscriptionID(subscription.ID, models.Upper)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting subscription with stripe ID %s\n", subscription.ID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		sub.Status = string(stripe.SubscriptionStatusCanceled)
		if err := models.Subscriptions.Update(*sub, models.Upper); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating subscription with stripe ID %s\n", subscription.ID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sub, err := models.Subscriptions.GetByStripeSubscriptionID(subscription.ID, models.Upper)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting subscription with stripe ID %s\n", subscription.ID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if subscription.CancelAtPeriodEnd {
			sub.Status = "active_until_period_end"
			if err := models.Subscriptions.Update(*sub, models.Upper); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating subscription with stripe ID %s\n", subscription.ID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			sub.Status = string(subscription.Status)
			if err := models.Subscriptions.Update(*sub, models.Upper); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating subscription with stripe ID %s\n", subscription.ID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	// ... handle other event types
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		w.WriteHeader(http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
