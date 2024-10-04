package internal

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
)

func NewHttpClientWithTLS(clientCertFile, clientKeyFile, caCertFile string) (http.Client, error) {
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return http.Client{}, err
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return http.Client{}, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			},
		},
	}, nil
}
