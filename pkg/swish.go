package pkg

import (
	"github.com/go-playground/validator/v10"
	"github.com/happsie/go-swish/internal"
)

type Swish struct {
	payment PaymentClient
}

func NewClient(clientCertFile, clientKeyFile, caCertFile string) (Swish, error) {
	client, err := internal.NewHttpClientWithTLS(clientCertFile, clientKeyFile, caCertFile)
	if err != nil {
		return Swish{}, err
	}
	return Swish{
		payment: PaymentClient{
			HttpClient: &client,
			Validator:  validator.New(validator.WithRequiredStructEnabled()),
			Host:       "mss.cpc.getswish.net",
			BaseURL:    "swish-cpcapi/api/v1",
		},
	}, nil
}

func (s *Swish) Payment() PaymentClient {
	return s.payment
}
