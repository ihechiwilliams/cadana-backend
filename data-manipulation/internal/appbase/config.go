package appbase

import (
	"time"

	"data-manipulation/pkg/config"
)

type Config struct {
	ServerAddress        string `env:"SERVER_ADDRESS" env-default:"0.0.0.0:3000"`
	ServiceName          string `env:"SERVICE_NAME" env-default:"cadana-backend"`
	ServerTimeout        int64  `env:"SERVER_TIMEOUT" env-default:"120"`
	Env                  string `env:"ENV"`
	LogLevel             string `env:"LOG_LEVEL" env-default:"debug"`
	SentryDSN            string `env:"SENTRY_DSN"`
	DdAgentHost          string `env:"DD_AGENT_HOST" env-default:"localhost"`
	CadanaBackendBaseURL string `env:"CADANA_BACKEND_BASE_URL" env-default:"localhost"`
	DdVersion            string `env:"DD_VERSION"`
}

func LoadConfig() (*Config, error) {
	c := new(Config)

	err := config.LoadConfig(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) HTTPServerTimeout() time.Duration {
	return time.Duration(c.ServerTimeout) * time.Second
}
