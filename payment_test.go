package goswish

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultRequest = PaymentRequest{
	PayeeAlias:            "1234679304",
	PayerAlias:            "4671234768",
	Amount:                1,
	Currency:              "SEK",
	CallbackURL:           "https://myfakehost.se/swishcallback.cfm",
	PayeePaymentReference: "0123456789",
	Message:               "Kingston USB Flash Drive 8 GB",
}

func TestCreatePaymentErrors(t *testing.T) {
	tests := []struct {
		statusCode    int
		body          PaymentRequest
		expectedError error
	}{
		{
			statusCode:    400,
			body:          defaultRequest,
			expectedError: BadRequestError,
		},
		{
			statusCode:    401,
			body:          defaultRequest,
			expectedError: UnauthorizedError,
		},
		{
			statusCode:    403,
			body:          defaultRequest,
			expectedError: ForbiddenError,
		},
		{
			statusCode:    415,
			body:          defaultRequest,
			expectedError: UnsupportedMediaTypeError,
		},
		{
			statusCode:    429,
			body:          defaultRequest,
			expectedError: TooManyRequestsError,
		},
		{
			statusCode:    500,
			body:          defaultRequest,
			expectedError: InternalServerError,
		},
		{
			statusCode:    422,
			body:          defaultRequest,
			expectedError: UnprocessableEntityError,
		},
	}

	for _, test := range tests {
		srv := newServer(t, test.statusCode, test.body)

		client, err := NewClient(Certificates{}, func(c *Config) {
			c.client = srv.Client()
			c.host = srv.URL
		})

		if err != nil {
			t.Fail()
		}

		_, _, err = client.Payment().Create(context.Background(), PaymentRequest{
			PayeeAlias:            "1234679304",
			PayerAlias:            "4671234768",
			Amount:                1,
			Currency:              "SEK",
			CallbackURL:           "https://myfakehost.se/swishcallback.cfm",
			PayeePaymentReference: "0123456789",
			Message:               "Kingston USB Flash Drive 8 GB",
		})
		if !errors.Is(err, test.expectedError) {
			t.Fail()
		}
		t.Logf("Test of error code %d: success", test.statusCode)
		srv.Close()
	}
}

func newServer(t *testing.T, statusCode int, body any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if body != nil {
			json, err := json.Marshal(body)
			if err != nil {
				t.Fail()
			}
			w.Write(json)
		}
	}))
}
