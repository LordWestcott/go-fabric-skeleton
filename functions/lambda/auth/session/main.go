package main

import (
	"appointme/data"
	"appointme/responses"
	"net/http"
	"time"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
	up "github.com/upper/db/v4"
)

type SessionResponse struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Token      string    `json:"token"`
	PictureURL string    `json:"picture_url"`
	Locale     string    `json:"locale"`
	StripeID   string    `json:"stripe_id"`
	CreatedAt  time.Time `json:"created_at"`
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

	sessionResponse := SessionResponse{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Token:      user.Token.PlainText,
		PictureURL: user.PictureURL,
		StripeID:   user.StripeID,
		Locale:     user.Locale,
		CreatedAt:  user.CreatedAt,
	}

	res := responses.NewResponseBase(true, "", sessionResponse)

	w.WriteHeader(http.StatusOK)
	res.WriteJson(w)
}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
