package cadana_backend

import (
	"context"
	cadanaBackendClient "data-manipulation/internal/clients/cadana_backend/openapi"
	"fmt"
	"net/http"
)

type Client interface {
	GetExchangeRate(ctx context.Context, params *cadanaBackendClient.V1GetExchangeRateJSONRequestBody) (*cadanaBackendClient.V1GetExchangeRateResponse, error)
}

type APIClient struct {
	client cadanaBackendClient.ClientWithResponsesInterface
}

func NewAPIClient(client cadanaBackendClient.ClientWithResponsesInterface) *APIClient {
	return &APIClient{
		client: client,
	}
}

func (a *APIClient) GetExchangeRate(ctx context.Context, params *cadanaBackendClient.V1GetExchangeRateJSONRequestBody) (*cadanaBackendClient.V1GetExchangeRateResponse, error) {
	response, err := a.client.V1GetExchangeRateWithResponse(ctx, *params)
	if err != nil {
		return nil, err
	}
	switch response.StatusCode() {
	case http.StatusOK:
		return response, nil
	default:
		return nil, a.apiError(response.StatusCode(), response.Body)
	}
}

func (a *APIClient) apiError(statusCode int, body []byte) error {
	return fmt.Errorf("unexpected error from Bureaucreat %d %s", statusCode, string(body))
}
