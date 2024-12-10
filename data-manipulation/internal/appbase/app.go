package appbase

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"

	"data-manipulation/internal/clients/cadana_backend"
	"data-manipulation/pkg/signals"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

type AppBase struct {
	Config              *Config
	DdStatsd            *statsd.Client
	Logger              zerolog.Logger
	ServiceName         string
	CadanaBackendClient *cadana_backend.APIClient
}

func New(options ...func(*AppBase)) *AppBase {
	appBase := &AppBase{}
	for _, o := range options {
		o(appBase)
	}

	return appBase
}

func Init(serviceName string) func(*AppBase) {
	return func(appBase *AppBase) {
		config, err := LoadConfig()
		if err != nil {
			panic(err)
		}

		appBase.Config = config
		appBase.ServiceName = serviceName
	}
}

func WithSignals(ctx context.Context, mainCtxStop context.CancelFunc) func(*AppBase) {
	return func(appBase *AppBase) {
		signals.HandleSignals(ctx, mainCtxStop, func() {})
	}
}

func WithSentry() func(*AppBase) {
	return func(appBase *AppBase) {
		fmt.Println(appBase.Config.SentryDSN)
		err := sentry.Init(sentry.ClientOptions{
			AttachStacktrace: true,
			Dsn:              appBase.Config.SentryDSN,
			Environment:      appBase.Config.Env,
			Release:          appBase.Config.DdVersion,
		})
		if err != nil {
			log.Logger.Error().Err(err).Msg("Sentry initialization failed")
		}
	}
}
