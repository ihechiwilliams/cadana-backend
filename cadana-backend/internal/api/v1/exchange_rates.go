package v1

import (
	"errors"
	"net/http"
	"sync"

	"cadana-backend/internal/api/server"
	"cadana-backend/internal/clients/servicea"
	"cadana-backend/internal/clients/serviceb"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/render"
)

type ExchangeRateHandler struct {
	serviceAClient servicea.ServiceAClient
	serviceBClient serviceb.ServiceBClient
}

func NewExchangeRateHandler(
	serviceAClient servicea.ServiceAClient,
	serviceBClient serviceb.ServiceBClient,
) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		serviceAClient: serviceAClient,
		serviceBClient: serviceBClient,
	}
}

func (a *API) V1GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	reqBody := new(server.V1GetExchangeRateJSONBody)

	err := render.DecodeJSON(r.Body, reqBody)
	if err != nil {
		sentry.CaptureException(err)
		server.BadRequestError(err, w, r)

		return
	}

	if reqBody.Data.CurrencyPair == "" {
		server.BadRequestError(errors.New("currency_pair is required"), w, r)
		return
	}

	exchangeRateData := reqBody.Data

	var wg sync.WaitGroup
	var mu sync.Mutex
	var rate float64

	done := make(chan struct{})

	wg.Add(2)

	// Call Service A
	go func() {
		defer wg.Done()

		if res, err := a.exchangeRateHandler.serviceAClient.ExchangeRate(r.Context(), servicea.ExchangeRateRequestBody{CurrencyPair: exchangeRateData.CurrencyPair}); err == nil {
			mu.Lock()
			select {
			case <-done:
			default:
				rate = res.Rate
				close(done)
			}
			mu.Unlock()
		}
	}()

	// Call Service B
	go func() {
		defer wg.Done()
		if res, err := a.exchangeRateHandler.serviceBClient.ExchangeRate(r.Context(), serviceb.ExchangeRateRequestBody{CurrencyPair: exchangeRateData.CurrencyPair}); err == nil {
			mu.Lock()
			select {
			case <-done:
			default:
				rate = res.Rate
				close(done)
			}
			mu.Unlock()
		}
	}()

	wg.Wait()

	select {
	case <-done:
		response := serializeRateToAPIResponse(exchangeRateData.CurrencyPair, rate)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	default:
		err = errors.New("unable to process request")
		server.ProcessingError(err, w, r)
		render.Status(r, http.StatusInternalServerError)
		return
	}
}

func serializeRateToAPIResponse(currencyPair string, rate float64) server.ExchangeRateResponse {
	return server.ExchangeRateResponse{
		Data: server.ExchangeRateResponseData{
			currencyPair: float32(rate),
		},
	}
}
