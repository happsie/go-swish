package goswish

import (
	"github.com/go-playground/validator/v10"
	"github.com/happsie/go-swish/internal"
)

type Swish struct {
	payment paymentClient
}

type Config struct {
	ClientCertFile string
	ClientKeyFile  string
	CaCertFile     string
	Host           string
}

func NewClient(config Config) (Swish, error) {
	client, err := internal.NewHttpClientWithTLS(config.ClientCertFile, config.ClientKeyFile, config.CaCertFile)
	if err != nil {
		return Swish{}, err
	}
	return Swish{
		payment: paymentClient{
			httpClient: &client,
			validator:  validator.New(validator.WithRequiredStructEnabled()),
			host:       config.Host,
			baseURL:    "swish-cpcapi/api/v1",
		},
	}, nil
}

func (s *Swish) Payment() paymentClient {
	return s.payment
}
