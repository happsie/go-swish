package main

import (
	"context"
	"fmt"

	goswish "github.com/happsie/go-swish"
)

// This example will teach you how to get the information regarding a payment from swish
// In order to retrieve a payment, it must first be created
func main() {
	client, _ := goswish.NewClient(goswish.Certificates{
		ClientCertFile: "/path/to/cert/public.pem",
		ClientKeyFile:  "/path/to/cert/private.key",
		CaCertFile:     "/path/to/cert/Swish_TLS_RootCA.pem",
	})

	location := "https://mss.cpc.getswish.net/swish-cpcapi/api/v1/paymentrequests/1CA5159969974F1F8CC9948A13FF643C"
	// GetInstructionID will try to parse the location received from Swish into only the ID
	instructionID, _ := goswish.GetInstructionID(location)

	ctx := context.Background()
	payment, _ := client.Payment().Retrieve(ctx, instructionID)
	
	fmt.Printf("payment retrieved: %v", payment)
}