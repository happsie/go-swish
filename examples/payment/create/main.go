package main

import (
	"context"
	"fmt"

	goswish "github.com/happsie/go-swish"
)

// This example will teach you how to create payment using swish
func main() {
	client, _ := goswish.NewClient(goswish.Certificates{
		ClientCertFile: "/path/to/cert/public.pem",
		ClientKeyFile:  "/path/to/cert/private.key",
		CaCertFile:     "/path/to/cert/Swish_TLS_RootCA.pem",
	})

	ctx := context.Background()
	location, paymentRequestToken, _ := client.Payment().Create(ctx, goswish.PaymentRequest{
		PayeeAlias:            "1234679304",
		PayerAlias:            "4671234768",
		Amount:                1,
		Currency:              "SEK",
		CallbackURL:           "https://myfakehost.se/swishcallback.cfm",
		PayeePaymentReference: "0123456789",
		Message:               "Kingston USB Flash Drive 8 GB",
	})
	fmt.Printf("Payment created: %s, %s", location, paymentRequestToken)
}
