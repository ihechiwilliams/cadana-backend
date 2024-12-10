package appbase

import (
	"cadana-backend/internal/api"
	v1 "cadana-backend/internal/api/v1"
	"cadana-backend/internal/clients/servicea"
	"cadana-backend/internal/clients/serviceb"
	"os"

	"cadana-backend/internal/api/server"
	openAPIUtils "cadana-backend/pkg/openapi"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/samber/do"
	"github.com/samber/lo"
)

func NewInjector(serviceName string, cfg *Config) *do.Injector {
	injector := do.New()

	// ===========================
	//	Service Configs (logging, open-api,...)
	// ===========================
	do.Provide(injector, func(i *do.Injector) (*zerolog.Logger, error) {
		logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			return nil, err
		}

		logger := zerolog.New(os.Stdout).
			Level(logLevel).
			With().
			Str("serviceName", serviceName).
			Logger()

		return &logger, nil
	})

	do.ProvideNamed(injector, InjectorApplicationRouter, func(i *do.Injector) (*chi.Mux, error) {
		logger := do.MustInvoke[*zerolog.Logger](i)
		openAPIValidation := do.MustInvokeNamed[*openAPIUtils.ValidationMiddleware](i, InjectorOpenAPIValidationMiddleware)

		return NewRouterMux(serviceName, logger, openAPIValidation, cfg.HTTPServerTimeout()), nil
	})

	do.ProvideNamed(injector, InjectorOpenAPIValidationMiddleware, func(i *do.Injector) (*openAPIUtils.ValidationMiddleware, error) {
		switch cfg.Env {
		case "test":
			return openAPIUtils.NewValidationMiddleware(
				openAPIUtils.WithDoc(lo.Must(server.GetSwagger())),
				openAPIUtils.WithErrorRenderer(server.ErrorRenderer),
			), nil
		default:
			return openAPIUtils.NewValidationMiddleware(
				openAPIUtils.WithDoc(lo.Must(server.GetSwagger())),
				openAPIUtils.WithKinOpenAPIDefaults(),
				openAPIUtils.WithErrorRenderer(server.ErrorRenderer),
			), nil
		}
	})

	// ===========================
	//	API services & Routes
	// ===========================
	do.Provide(injector, func(i *do.Injector) (*v1.ExchangeRateHandler, error) {
		return v1.NewExchangeRateHandler(
			do.MustInvoke[servicea.ServiceAClient](i),
			do.MustInvoke[serviceb.ServiceBClient](i),
		), nil
	})

	do.Provide(injector, func(i *do.Injector) (*v1.API, error) {
		exchangeRateHandler := do.MustInvoke[*v1.ExchangeRateHandler](i)

		return v1.NewAPI(exchangeRateHandler), nil
	})

	do.Provide(injector, func(i *do.Injector) (*api.Routes, error) {
		v1API := do.MustInvoke[*v1.API](i)

		return api.NewRoutes(v1API), nil
	})

	// ===========================
	//	External Services & Routes
	// ===========================
	do.Provide(injector, func(i *do.Injector) (servicea.ServiceAClient, error) {
		return servicea.NewClient(
			cfg.ServiceAAPIURL,
			cfg.ServiceAAPIKey,
		), nil
	})

	do.Provide(injector, func(i *do.Injector) (serviceb.ServiceBClient, error) {
		return serviceb.NewClient(
			cfg.ServiceAAPIURL,
			cfg.ServiceAAPIKey,
		), nil
	})

	return injector
}
