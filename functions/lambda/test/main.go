package main

import (
	"fmt"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
)

func handler(w http.ResponseWriter, r *http.Request) {

	app, err := gofabric.InitApp()
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	message := ""
	if app.DB != nil {
		message += "DB is connected\n"
	} else {
		message += "DB is not connected\n"
	}

	if app.Google_OAuth2 != nil {
		message += "GoogleSignIn is connected\n"
	} else {
		message += "GoogleSignIn is not connected\n"
	}

	if app.Messaging != nil {
		message += "Messaging is connected\n"

		if app.Messaging.EmailService != nil {
			message += "EmailService is connected\n"
		} else {
			message += "EmailService is not connected\n"
		}

		if app.Messaging.SMSService != nil {
			message += "SMS is connected\n"
		} else {
			message += "SMS is not connected\n"
		}

		if app.Messaging.VerificationService != nil {
			message += "VerificationService is connected\n"
		} else {
			message += "VerificationService is not connected\n"
		}

		if app.Messaging.WhatsAppService != nil {
			message += "WhatsAppService is connected\n"
		} else {
			message += "WhatsAppService is not connected\n"
		}

	} else {
		message += "Messaging is not connected\n"
	}

	if app.OpenAI != nil {
		message += "OpenAI is connected\n"
	} else {
		message += "OpenAI is not connected\n"
	}

	fmt.Println("got here 9")

	if app.Stripe != nil {
		message += "Stripe is connected\n"
	} else {
		message += "Stripe is not connected\n"
	}

	fmt.Fprintf(w, message)
}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
