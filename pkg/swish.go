package pkg

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
)

type SwishClient struct {
	client    http.Client
	validator *validator.Validate
	host      string
	baseURL   string
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

func NewClient(clientCertFile, clientKeyFile, caCertFile string) (SwishClient, error) {
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return SwishClient{}, err
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return SwishClient{}, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			},
		},
	}
	return SwishClient{
		client:    client,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		host:      "mss.cpc.getswish.net",
		baseURL:   "swish-cpcapi/api/v1",
	}, nil
}

func (sc *SwishClient) CreatePayment(ctx context.Context, request PaymentRequest) (string, error) {
	err := sc.validator.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return "", fmt.Errorf("validation of request failed. Field %s is required", err.StructField())
		}
	}
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://%s/%s/paymentrequests", sc.host, sc.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := sc.client.Do(req)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return res.Header.Get("Location"), nil
}
