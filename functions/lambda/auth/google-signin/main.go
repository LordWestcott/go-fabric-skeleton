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

	app.Google_OAuth2.SignIn(w, r)
}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
