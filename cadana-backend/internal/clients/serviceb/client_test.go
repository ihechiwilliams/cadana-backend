package serviceb

import (
	"context"
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ServiceBClientSuite struct {
	suite.Suite
	client  *Client
	metrics statsd.ClientInterface
}

func (s *ServiceBClientSuite) SetupTest() {
	s.metrics, _ = statsd.New("localhost")
	s.client = NewClient("https://api.serviceb.com", "test-secret-key")

	httpmock.Activate()
}

func (s *ServiceBClientSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (s *ServiceBClientSuite) TestServiceBClient_ExchangeRate() {
	s.Run("get exchange rate is successful", func() {
		httpmock.RegisterResponder(
			http.MethodPost,
			"https://api.serviceb.com/v1/exchange_rate",
			httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("testdata/exchange_rate_success_response.json")),
		)

		response, err := s.client.ExchangeRate(context.Background(), ExchangeRateRequestBody{
			CurrencyPair: "USD/EUR",
		})
		s.NoError(err)
		s.Equal(10.00, response.Rate)
	})

	s.Run("get exchange rate returns error", func() {
		httpmock.RegisterResponder(
			http.MethodPost,
			"https://api.serviceb.com/v1/exchange_rate",
			httpmock.NewJsonResponderOrPanic(http.StatusNotFound, httpmock.File("testdata/exchange_rate_error_response.json")),
		)

		response, err := s.client.ExchangeRate(context.Background(), ExchangeRateRequestBody{
			CurrencyPair: "USD/EUR",
		})

		s.Nil(response)
		s.Error(err)
	})
}

func TestServiceBClientSuite(t *testing.T) {
	suite.Run(t, new(ServiceBClientSuite))
}
