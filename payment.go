package goswish

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type paymentClient struct {
	conf Config
}

type PaymentRequest struct {
	PayeeAlias            string  `json:"payeeAlias" validate:"required"`
	PayerAlias            string  `json:"payerAlias"`
	Amount                float32 `json:"amount" validate:"required"`
	Currency              string  `json:"currency" validate:"required"`
	CallbackURL           string  `json:"callbackUrl" validate:"required"`
	PayeePaymentReference string  `json:"payeePaymentReference"`
	Message               string  `json:"message"`
}

type Payment struct {
	PayeeAlias            string  `json:"payeeAlias" validate:"required"`
	PayerAlias            string  `json:"payerAlias"`
	Amount                float32 `json:"amount" validate:"required"`
	Currency              string  `json:"currency" validate:"required"`
	CallbackURL           string  `json:"callbackUrl" validate:"required"`
	PayeePaymentReference string  `json:"payeePaymentReference"`
	Message               string  `json:"message"`
	ID                    string  `json:"id"`
	Status                string  `json:"status"`
	DateCreated           string  `json:"dateCreated"`
	DatePaid              string  `json:"datePaid"`
}

type cancelPayment struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

const (
	contentTypeJson      = "application/json"
	contentTypePatchJson = "application/json-patch+json"
)

// Create creates a payment with information provided in the PaymentRequest.
// The HTTP call is created with the provided context.
// Returns the location, paymentRequestToken (Only for M-Commerce), or an error.
func (sc paymentClient) Create(ctx context.Context, request PaymentRequest) (location string, paymentRequestToken string, err error) {
	err = sc.conf.validator.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return "", "", fmt.Errorf("validation of request failed. Field %s is required", err.StructField())
		}
	}
	body, err := json.Marshal(request)
	if err != nil {
		return "", "", err
	}
	url := fmt.Sprintf("%s/%s/paymentrequests", sc.conf.host, sc.conf.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", contentTypeJson)
	req.Header.Set("Accept", contentTypeJson)
	res, err := sc.conf.client.Do(req)
	if err != nil {
		return "", "", err
	}
	return res.Header.Get("Location"), res.Header.Get("PaymentRequestToken"), nil
}

// Retrieve retrieves a payment with information created from the Create method.
// The HTTP call is created with the provided context.
// Use GetInstructionID util function to extract the ID from Location from the Create method
// Returns the payment or an error
func (sc paymentClient) Retrieve(ctx context.Context, ID string) (Payment, error) {
	url := fmt.Sprintf("%s/%s/paymentrequests/%s", sc.conf.host, sc.conf.baseURL, ID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Payment{}, err
	}
	req.Header.Set("Content-Type", contentTypeJson)
	req.Header.Set("Accept", contentTypeJson)

	res, err := sc.conf.client.Do(req)
	if err != nil {
		return Payment{}, err
	}
	defer res.Body.Close()

	var payment Payment
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&payment)
	if err != nil {
		return Payment{}, err
	}

	return payment, nil
}

// Cancel cancels a payment identified by the instructionID retrieved from the Create method.
// The HTTP call is created with the provided context.
// Use GetInstructionID util function to extract the ID from Location from the Create method
// Returns the payment or an error
func (sc paymentClient) Cancel(ctx context.Context, ID string) (Payment, error) {
	url := fmt.Sprintf("%s/%s/paymentrequests/%s", sc.conf.host, sc.conf.baseURL, ID)

	cancel := []cancelPayment{{
		Op:    "replace",
		Path:  "/status",
		Value: "cancelled",
	}}
	body, err := json.Marshal(cancel)
	if err != nil {
		return Payment{}, err
	}
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(body))
	if err != nil {
		return Payment{}, err
	}
	req.Header.Set("Content-Type", contentTypePatchJson)
	req.Header.Set("Accept", contentTypeJson)

	res, err := sc.conf.client.Do(req)
	if err != nil {
		return Payment{}, err
	}
	defer res.Body.Close()

	var payment Payment
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&payment)
	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}
