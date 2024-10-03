package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type PaymentClient struct {
	HttpClient *http.Client
	Validator  *validator.Validate
	Host       string
	BaseURL    string
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

func (sc PaymentClient) Create(ctx context.Context, request PaymentRequest) (string, error) {
	err := sc.Validator.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return "", fmt.Errorf("validation of request failed. Field %s is required", err.StructField())
		}
	}
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://%s/%s/paymentrequests", sc.Host, sc.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := sc.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return res.Header.Get("Location"), nil
}
