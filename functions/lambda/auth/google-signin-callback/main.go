package main

import (
	"appointme/data"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
	"github.com/lordwestcott/gofabric/jwt"
	"github.com/stripe/stripe-go/v74"

	up "github.com/upper/db/v4"
)

type Google_Account struct {
	ID             string `json:"id"`
	FirstName      string `json:"given_name"`
	LastName       string `json:"family_name"`
	PictureUrl     string `json:"picture"`
	Locale         string `json:"locale"`
	Email          string `json:"email"`
	Verified_Email bool   `json:"verified_email"`
}

func (ga *Google_Account) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, ga); err != nil {
		return err
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	app, err := gofabric.InitApp()
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	models := data.New(app.DB)

	raw, err := app.Google_OAuth2.CallBack(w, r)
	if err != nil {
		http.Redirect(w, r, os.Getenv("REDIRECT_SOMETHING_WENT_WRONG_PAGE"), http.StatusTemporaryRedirect)
	}

	// marshal the data into a struct
	ga := Google_Account{}
	ga.Unmarshal(raw)

	// get user from the database or create a new one
	user, err := models.Users.GetByGoogleID(ga.ID, models.Upper)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			app.ErrorLog.Println("USER GET BY GOOGLE ID ERROR:", err)
			return
		}
		// create a new user
		user = &data.User{
			FirstName:       ga.FirstName,
			LastName:        ga.LastName,
			Email:           ga.Email,
			Active:          1,
			GoogleID:        ga.ID,
			Password:        "",
			EmailIsVerified: ga.Verified_Email,
			PictureURL:      ga.PictureUrl,
			Locale:          ga.Locale,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		_, err = models.Users.Insert(user, models.Upper)
		if err != nil {
			app.ErrorLog.Println("USER INSERT GOOGLE USER ERROR:", err)
			return
		}
		fmt.Println("USER INSERTED: ID -> ", user.ID)

		if app.Stripe != nil {
			fmt.Println("SETTING UP STRIPE CUSTOMER")
			// Create a new user on Stripe
			params := &stripe.CustomerParams{
				Email: stripe.String(user.Email),
				Name:  stripe.String(user.FirstName + " " + user.LastName),
				Params: stripe.Params{
					Metadata: map[string]string{
						"uID": strconv.FormatInt(user.ID, 10),
					},
				},
			}
			customer, err := app.Stripe.CreateCustomer(params)
			if err != nil {
				app.ErrorLog.Println("STRIPE CREATE CUSTOMER ERROR:", err)
				return
			}

			// Update the user with the Stripe ID
			user.StripeID = customer.ID
			fmt.Println("UPDATING USER WITH STRIPE ID:", user.StripeID)
			err = models.Users.Update(*user, models.Upper)
			if err != nil {
				app.ErrorLog.Println("USER UPDATE STRIPE ID ERROR:", err)
				return
			}
		}
	}

	claims := jwt.Claims{
		UserID:   user.ID,
		Username: user.FirstName + " " + user.LastName,
		Email:    user.Email,
		Scope:    []string{"user"},
	}

	jwt, err := app.JWT.GenerateJWT(&claims, int(time.Duration(24*time.Hour).Seconds()))
	if err != nil {
		app.ErrorLog.Println("Generate JWT ERROR:", err)
	}

	token := data.Token{}
	t, err := token.GenerateToken(user.ID, jwt, time.Duration(24*time.Hour))
	if err != nil {
		app.ErrorLog.Println("Generate Token ERROR:", err)
	}
	token.Insert(*t, *user, models.Upper)

	redirect := os.Getenv("FRONTEND_URL") + os.Getenv("REDIRECT_LOGIN_SUCCESS") + "?token=" + t.PlainText
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
