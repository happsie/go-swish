package pkg

import (
	"github.com/go-playground/validator/v10"
	"github.com/happsie/go-swish/internal"
)

type Swish struct {
	payment paymentClient
}

func NewClient(clientCertFile, clientKeyFile, caCertFile string) (Swish, error) {
	client, err := internal.NewHttpClientWithTLS(clientCertFile, clientKeyFile, caCertFile)
	if err != nil {
		return Swish{}, err
	}
	return Swish{
		payment: paymentClient{
			httpClient: &client,
			validator:  validator.New(validator.WithRequiredStructEnabled()),
			host:       "https://mss.cpc.getswish.net",
			baseURL:    "swish-cpcapi/api/v1",
		},
	}, nil
}

func (s *Swish) Payment() paymentClient {
	return s.payment
}
