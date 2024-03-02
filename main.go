package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/ephemeralkey"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

func main() {
	// This is your test secret API key.
	stripe.Key = "sk_test_51OpcPfHVyvJVrH9E5AEuUFa7ynSwVCL2KhaI6xtOSTfmSjBFRUnn2O6tDlYOkyjHsfawONwvb6LPE4JCbnw2g90k0011fynIOm"

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/create-payment-intent", handlePaymentSheet)

	addr := "localhost:8081"
	log.Printf("Listening on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handlePaymentSheet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Use an existing Customer ID if this is a returning customer.
	cparams := &stripe.CustomerParams{}
	c, _ := customer.New(cparams)

	ekparams := &stripe.EphemeralKeyParams{
		Customer:      stripe.String(c.ID),
		StripeVersion: stripe.String("2023-10-16"),
	}
	ek, _ := ephemeralkey.New(ekparams)

	piparams := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(1099),
		Currency: stripe.String(string(stripe.CurrencyEUR)),
		Customer: stripe.String(c.ID),
		PaymentMethodTypes: []*string{
			stripe.String("bancontact"),
			stripe.String("card"),
			stripe.String("ideal"),
			stripe.String("klarna"),
			stripe.String("sepa_debit"),
		},
	}
	pi, _ := paymentintent.New(piparams)

	writeJSON(w, struct {
		PaymentIntent  string `json:"paymentIntent"`
		EphemeralKey   string `json:"ephemeralKey"`
		Customer       string `json:"customer"`
		PublishableKey string `json:"publishableKey"`
	}{
		PaymentIntent:  pi.ClientSecret,
		EphemeralKey:   ek.Secret,
		Customer:       c.ID,
		PublishableKey: "pk_test_51OpcPfHVyvJVrH9E6fZD7MehGxcmrvAGaX1y1gWrQrUgXlki94f9MlnxWUoheeE6my1fITxokmAojqHtjQisolxp00VpgM8Zb3",
	})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
