package main

import (
	"appointme/data"
	"appointme/responses"
	"fmt"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/lordwestcott/gofabric"
)

//This has been set up specifically for subscriptions.
//But will work for any product with prices.

type Response struct {
	Products []ProductDTO `json:"products"`
}

type ProductDTO struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsRecurring bool       `json:"is_recurring"`
	Prices      []PriceDTO `json:"prices"`
}

type PriceDTO struct {
	ID            int64  `json:"id"`
	Amount        int    `json:"amount"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	BillingPeriod string `json:"billing_period"`
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

	pd := data.Product{}
	pr := data.Price{}

	products, err := pd.GetAll(models.Upper)
	if err != nil {
		res := responses.NewResponseBase(false, err.Error(), nil)
		res.WriteJson(w)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := Response{}

	for _, product := range products {
		prices, err := pr.GetByProductID(product.ID, models.Upper)
		if err != nil {
			res := responses.NewResponseBase(false, err.Error(), nil)
			res.WriteJson(w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		priceDTOarr := []PriceDTO{}
		for _, price := range prices {
			fmt.Printf("Price: %v\n", price.ID)
			pdto := PriceDTO{
				ID:            price.ID,
				Amount:        price.Amount,
				Name:          price.Name,
				Description:   price.Description,
				BillingPeriod: price.BillingPeriod,
			}
			priceDTOarr = append(priceDTOarr, pdto)
		}

		productDto := ProductDTO{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			IsRecurring: product.IsRecurring,
			Prices:      priceDTOarr,
		}

		res.Products = append(res.Products, productDto)
	}

	resBase := responses.NewResponseBase(true, "", res)
	resBase.WriteJson(w)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
