# Create Payment

Below is an explanation of the example with each line described with a comment

```go
func main() {
    // Setup the goswish client, this is the client you will use to integrate with swish
    // Important: It is possible to overload the HttpClient in case you don't want to provide certificates like this.
    // Overloading of options can be sent as a second argument to the NewClient function.
	client, _ := goswish.NewClient(goswish.Certificates{
        // Specify the filepath to ClientCertFile
		ClientCertFile: "/path/to/cert/public.pem",
        // Specify the filepath to ClientKeyFile
		ClientKeyFile:  "/path/to/cert/private.key",
        // Specify the filepath to the Root CA
		CaCertFile:     "/path/to/cert/Swish_TLS_RootCA.pem",
	})

    // Create a context
	ctx := context.Background()
    // Call the Payment.Create to create a payment. Returns the location, PaymentRequestToken (M-Commerce) and an error
	location, paymentRequestToken, _ := client.Payment().Create(ctx, goswish.PaymentRequest{
		PayeeAlias:            "1234679304",
		PayerAlias:            "4671234768",
		Amount:                1,
		Currency:              "SEK",
		CallbackURL:           "https://myfakehost.se/swishcallback.cfm",
		PayeePaymentReference: "0123456789",
		Message:               "Kingston USB Flash Drive 8 GB",
	})
    // Print the created payment to stdout
	fmt.Printf("Payment created: %s, %s", location, paymentRequestToken)
}
```