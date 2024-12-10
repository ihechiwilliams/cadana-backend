package servicea

import (
	"context"
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ServiceAClientSuite struct {
	suite.Suite
	client  *Client
	metrics statsd.ClientInterface
}

func (s *ServiceAClientSuite) SetupTest() {
	s.metrics, _ = statsd.New("localhost")
	s.client = NewClient("https://api.servicea.com", "test-secret-key")

	httpmock.Activate()
}

func (s *ServiceAClientSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (s *ServiceAClientSuite) TestServiceAClient_ExchangeRate() {
	s.Run("get exchange rate is successful", func() {
		httpmock.RegisterResponder(
			http.MethodPost,
			"https://api.servicea.com/v1/exchange_rate",
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
			"https://api.servicea.com/v1/exchange_rate",
			httpmock.NewJsonResponderOrPanic(http.StatusNotFound, httpmock.File("testdata/exchange_rate_error_response.json")),
		)

		response, err := s.client.ExchangeRate(context.Background(), ExchangeRateRequestBody{
			CurrencyPair: "USD/EUR",
		})

		s.Nil(response)
		s.Error(err)
	})
}

func TestServiceAClientSuite(t *testing.T) {
	suite.Run(t, new(ServiceAClientSuite))
}
