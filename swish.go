package goswish

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/happsie/go-swish/internal"
)

type Swish struct {
	payment paymentClient
}

type Config struct {
	host      string
	client    *http.Client
	baseURL   string
	validator *validator.Validate
}

type Certificates struct {
	ClientCertFile string
	ClientKeyFile  string
	CaCertFile     string
}

type Option func(*Config)

func NewClient(cert Certificates, options ...Option) (Swish, error) {
	c := Config{
		host:      "https://mss.cpc.getswish.net",
		baseURL:   "swish-cpcapi/api/v1",
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	for _, opt := range options {
		opt(&c)
	}

	// Setup default if no httpClient was provided by options
	if c.client == nil {
		client, err := internal.NewHttpClientWithTLS(cert.ClientCertFile, cert.ClientKeyFile, cert.CaCertFile)
		if err != nil {
			return Swish{}, err
		}
		c.client = &client
	}

	return Swish{
		payment: paymentClient{
			conf: c,
		},
	}, nil
}

func (s *Swish) Payment() paymentClient {
	return s.payment
}
