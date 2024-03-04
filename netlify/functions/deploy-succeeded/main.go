package main

import (
	"encoding/json"
	"log"
	"os"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/ephemeralkey"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		log.Print("Invalid HTTP method: ", request.HTTPMethod)
		return &events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusMethodNotAllowed), StatusCode: http.StatusMethodNotAllowed}, nil
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

	responseBody, _ := json.Marshal(struct {
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

	return &events.APIGatewayProxyResponse{Body: string(responseBody), StatusCode: http.StatusOK}, nil
}

func main() {
	// This is your test secret API key.
	stripe.Key = "sk_test_51OpcPfHVyvJVrH9E5AEuUFa7ynSwVCL2KhaI6xtOSTfmSjBFRUnn2O6tDlYOkyjHsfawONwvb6LPE4JCbnw2g90k0011fynIOm"

	// Define AWS Lambda environment variables
	os.Setenv("_LAMBDA_SERVER_PORT", "3000")
	//os.Setenv("AWS_LAMBDA_RUNTIME_API", "runtime.somewhere.com")
	//lambda.StartHandler(&handler{})
	lambda.Start(handler)
}
