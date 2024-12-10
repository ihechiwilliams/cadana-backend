package appbase

import (
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/samber/oops"
)

type StatsdService struct {
	Client statsd.ClientInterface
}

func NewStatsdService(ddAgentHost string) (*StatsdService, error) {
	client, err := statsd.New(ddAgentHost)
	if err != nil {
		return nil, err
	}

	return &StatsdService{Client: client}, nil
}

func (s *StatsdService) HealthCheck() error {
	if s.Client.IsClosed() {
		return oops.Errorf("client is closed")
	}

	return nil
}

func (s *StatsdService) Shutdown() error {
	err := s.Client.Flush()
	if err != nil {
		return err
	}

	err = s.Client.Close()
	if err != nil {
		return err
	}

	return nil
}
