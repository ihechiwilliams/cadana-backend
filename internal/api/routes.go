package api

import (
	"cadana-backend/internal/api/server"
	v1 "cadana-backend/internal/api/v1"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(router *chi.Mux, si *Routes) {
	server.HandlerFromMux(si, router)
}

func NewRoutes(apiV1 *v1.API) *Routes {
	return &Routes{
		v1: apiV1,
	}
}

// Routes is the wrapper for all the versions of the API defined by server.ServerInterface.
type Routes struct {
	v1 *v1.API
}

func (a Routes) V1GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	a.v1.V1GetExchangeRates(w, r)
}
