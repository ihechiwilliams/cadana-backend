package v1

type API struct {
	exchangeRateHandler *ExchangeRateHandler
}

func NewAPI(exchangeRateHandler *ExchangeRateHandler) *API {
	return &API{
		exchangeRateHandler: exchangeRateHandler,
	}
}
