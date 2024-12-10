package v1_test

import (
	"errors"
	"net/http"
	"testing"

	"cadana-backend/internal/api"
	v1 "cadana-backend/internal/api/v1"
	"cadana-backend/internal/clients/servicea"
	serviceAMock "cadana-backend/internal/clients/servicea/mocks"
	"cadana-backend/internal/clients/serviceb"
	serviceBMock "cadana-backend/internal/clients/serviceb/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type exchangeRatesSuite struct {
	suite.Suite
	mux      *chi.Mux
	serviceA *serviceAMock.MockServiceAClient
	serviceB *serviceBMock.MockServiceBClient
}

func (e *exchangeRatesSuite) SetupTest() {
	e.mux = chi.NewMux()

	//e.serviceA = &serviceAMock.MockServiceAClient{}
	//e.serviceA = &serviceAMock.MockServiceAClient{}

	e.serviceA = serviceAMock.NewMockServiceAClient(e.T())
	e.serviceB = serviceBMock.NewMockServiceBClient(e.T())

	exchangeRatesHandler := v1.NewExchangeRateHandler(e.serviceA, e.serviceB)
	apiService := v1.NewAPI(exchangeRatesHandler)
	api.InitRoutes(e.mux, api.NewRoutes(apiService))

}

func TestExchangeRate(t *testing.T) {
	suite.Run(t, new(exchangeRatesSuite))
}

func (e *exchangeRatesSuite) Test_V1GetExchangeRates() {
	path := "/v1/exchange-rate"

	e.Run("when Service A returns a successful response", func() {
		requestBody := `{
			"data": {
				"currency_pair": "USD-EUR"
			}
		}`

		e.serviceA.
			On(
				"ExchangeRate",
				mock.Anything,
				servicea.ExchangeRateRequestBody{
					CurrencyPair: "USD-EUR",
				},
			).
			Return(&servicea.ExchangeRateResponse{Rate: 1.1}, nil).
			Once()

		e.serviceB.
			On(
				"ExchangeRate",
				mock.Anything,
				serviceb.ExchangeRateRequestBody{
					CurrencyPair: "USD-EUR",
				},
			).
			Return(nil, errors.New("service A error")).
			Once()

		apitest.New().
			Handler(e.mux).
			Post(path).
			Header("Content-Type", "application/json").
			Body(requestBody).
			Expect(e.T()).
			Status(http.StatusOK).
			Body(`{"data":{"USD-EUR":1.1}}`).
			End()
	})

	e.Run("when Service B returns a successful response after Service A fails", func() {
		requestBody := `{
			"data": {
				"currency_pair": "USD-EUR"
			}
		}`

		e.serviceA.
			On(
				"ExchangeRate",
				mock.Anything,
				servicea.ExchangeRateRequestBody{
					CurrencyPair: "USD-EUR",
				},
			).
			Return(nil, errors.New("service A error")).
			Once()

		e.serviceB.
			On(
				"ExchangeRate",
				mock.Anything,
				serviceb.ExchangeRateRequestBody{
					CurrencyPair: "USD-EUR",
				},
			).
			Return(&serviceb.ExchangeRateResponse{Rate: 10.0}, nil).
			Once()

		apitest.New().
			Handler(e.mux).
			Post(path).
			Header("Content-Type", "application/json").
			Body(requestBody).
			Expect(e.T()).
			Status(http.StatusOK).
			Body(`{"data":{"USD-EUR":10.0}}`).
			End()
	})

	e.Run("when both services fail", func() {
		requestBody := `{
			"data": {
				"currency_pair": "USD-EUR"
			}
		}`

		e.serviceA.
			On(
				"ExchangeRate",
				mock.Anything,
				servicea.ExchangeRateRequestBody{
					CurrencyPair: "USD-EUR",
				},
			).
			Return(nil, errors.New("service A error")).
			Once()

		e.serviceB.
			On(
				"ExchangeRate",
				mock.Anything,
				serviceb.ExchangeRateRequestBody{
					CurrencyPair: "USD-EUR",
				},
			).
			Return(nil, errors.New("service B error")).
			Once()

		apitest.New().
			Handler(e.mux).
			Post(path).
			Header("Content-Type", "application/json").
			Body(requestBody).
			Expect(e.T()).
			Status(http.StatusUnprocessableEntity).
			Body(`{
				"errors": [
					{
						"code": "PROCESSING_ERROR",
						"detail": "unable to process request",
						"meta": {},
						"status": 422,
						"title": "PROCESSING_ERROR"
					}
				]
			}`).
			End()
	})

	e.Run("when invalid request body is provided", func() {
		invalidRequestBody := `{
			"data": {}
		}`

		apitest.New().
			Handler(e.mux).
			Post(path).
			Header("Content-Type", "application/json").
			Body(invalidRequestBody).
			Expect(e.T()).
			Status(http.StatusBadRequest).
			Body(`{
				"errors": [
					{
						"code": "Bad Request",
						"detail": "currency_pair is required",
						"meta": {},
						"status": 400,
						"title": "BAD_REQUEST"
					}
				]
			}`).
			End()
	})
}
