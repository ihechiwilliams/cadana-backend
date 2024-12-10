package appbase

import (
	"time"

	"cadana-backend/pkg/config"
)

type Config struct {
	ServerAddress  string `env:"SERVER_ADDRESS" env-default:"0.0.0.0:3000"`
	ServiceName    string `env:"SERVICE_NAME" env-default:"cadana-backend"`
	Port           string `env:"PORT" env-default:"3000"`
	ServerTimeout  int64  `env:"SERVER_TIMEOUT" env-default:"120"`
	Env            string `env:"ENV"`
	LogLevel       string `env:"LOG_LEVEL" env-default:"debug"`
	SentryDSN      string `env:"SENTRY_DSN"`
	DdAgentHost    string `env:"DD_AGENT_HOST" env-default:"localhost"`
	ServiceAAPIURL string `env:"SERVICEA_API_URL" env-required:"true"`
	ServiceAAPIKey string `env:"SERVICEA_API_KEY" env-required:"true"`
	ServiceBAPIURL string `env:"SERVICEB_API_URL" env-required:"true"`
	ServiceBAPIKey string `env:"SERVICEB_API_KEY" env-required:"true"`
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
