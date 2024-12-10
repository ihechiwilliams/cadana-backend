package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ExchangeRates struct {
	apiURL string
	cache  map[string]float64 // Cache for exchange rates to minimize API calls
}

// NewExchangeRates initializes the ExchangeRates with an API URL and empty cache
func NewExchangeRates(apiURL string) *ExchangeRates {
	return &ExchangeRates{
		apiURL: apiURL,
		cache:  make(map[string]float64),
	}
}

// GetRate fetches the exchange rate between two currencies
func (er *ExchangeRates) GetRate(from, to string) float64 {
	currencyPair := fmt.Sprintf("%s-%s", from, to)
	// Check cache first
	if rate, found := er.cache[currencyPair]; found {
		return rate
	}

	// Fetch rate from external API
	rate, err := er.fetchExchangeRate(currencyPair)
	if err != nil {
		fmt.Printf("Error fetching exchange rate for %s: %v\n", currencyPair, err)
		return 0.0 // Return 0.0 on error to indicate failure
	}

	// Cache the fetched rate
	er.cache[currencyPair] = rate
	return rate
}

// fetchExchangeRate makes a POST request to fetch the exchange rate
func (er *ExchangeRates) fetchExchangeRate(currencyPair string) (float64, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"data": map[string]string{
			"currency_pair": currencyPair,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return 0.0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make the HTTP POST request
	resp, err := http.Post(er.apiURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return 0.0, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return 0.0, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	// Decode the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0.0, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract the exchange rate
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return 0.0, fmt.Errorf("invalid response structure")
	}
	rate, ok := data[currencyPair].(float64)
	if !ok {
		return 0.0, fmt.Errorf("exchange rate not found in response")
	}

	return rate, nil
}
