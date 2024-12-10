package servicea

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

const (
	exchangeRatePath = "v1/exchange_rate"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	secretKey  string
}

type ServiceAClient interface {
	ExchangeRate(ctx context.Context, requestBody ExchangeRateRequestBody) (*ExchangeRateResponse, error)
}

func NewClient(
	baseURL string,
	secretKey string,
) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		secretKey:  secretKey,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Probably Auth Required
	req.Header.Add("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *Client) ExchangeRate(ctx context.Context, requestBody ExchangeRateRequestBody) (*ExchangeRateResponse, error) {
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fullURL := fmt.Sprintf("%s/%s", c.baseURL, exchangeRatePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("api-key", c.secretKey)

	response, err := c.Do(req)
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var exchangeRateResponse ExchangeRateResponse
	respErr := checkResponse(response, &exchangeRateResponse)
	if respErr != nil {
		return nil, errors.WithStack(respErr)
	}

	return &exchangeRateResponse, nil
}

func checkResponse(res *http.Response, successBody interface{}) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("[service_a] API error, could not read the response body %v", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("[service_a] API error, message: %v: status code %d", string(body), res.StatusCode)
	} else { // success case
		err = json.Unmarshal(body, &successBody)
		if err != nil {
			return fmt.Errorf("[service_a] API error, success case message is not decode into %T: err: %v", successBody, err)
		}

		return nil
	}
}
