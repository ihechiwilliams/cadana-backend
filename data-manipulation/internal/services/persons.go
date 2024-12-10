package services

import (
	"context"
	"data-manipulation/internal/clients/cadana_backend"
	"data-manipulation/internal/clients/cadana_backend/openapi"
	"data-manipulation/internal/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"sort"
)

type ServicePersons struct {
	Data []models.Person
}

// SortBySalaryAscending Sort the data array of Person objects by salary in ascending order.
func (ps *ServicePersons) SortBySalaryAscending() {
	sort.Slice(ps.Data, func(i, j int) bool {
		return ps.Data[i].GetSalaryAsFloat() < ps.Data[j].GetSalaryAsFloat()
	})
}

// SortBySalaryDescending Sort the data array of Person objects by salary in descending order.
func (ps *ServicePersons) SortBySalaryDescending() {
	sort.Slice(ps.Data, func(i, j int) bool {
		return ps.Data[i].GetSalaryAsFloat() > ps.Data[j].GetSalaryAsFloat()
	})
}

// GroupByCurrency Group Person objects by salary currency into a hash map.
func (ps *ServicePersons) GroupByCurrency() map[string][]models.Person {
	grouped := make(map[string][]models.Person)
	for _, person := range ps.Data {
		currency := person.Salary.Currency
		grouped[currency] = append(grouped[currency], person)
	}
	return grouped
}

// FilterBySalaryInUSD Filter Person objects by salary in USD based on a threshold.
// Uses the ExchangeRateProvider interface for currency conversion.
func (ps *ServicePersons) FilterBySalaryInUSD(threshold float64, exchangeRates cadana_backend.APIClient) []models.Person {
	var filtered []models.Person
	for _, person := range ps.Data {
		var exRate float32
		res, err := exchangeRates.GetExchangeRate(context.Background(), &openapi.V1GetExchangeRateJSONRequestBody{
			Data: &openapi.ExchangeRateRequestBody{
				CurrencyPair: person.Salary.Currency + "-" + "USD",
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch exchange rate")

			return filtered
		}

		// Check if the response is successful and contains the expected data
		if res.JSON200 != nil {
			fmt.Println("here")
			// Access the exchange rate data
			data := res.JSON200.Data

			// Retrieve the value using the currency pair as the key
			currencyPair := person.Salary.Currency + "-" + "USD"
			exRate = data[currencyPair]
		} else {
			fmt.Println("here1")
			// Handle error responses (e.g., 400, 422, 500)
			if res.JSON400 != nil {
				fmt.Printf("Error 400")
			} else if res.JSON422 != nil {
				fmt.Printf("Error 422")
			} else if res.JSON500 != nil {
				fmt.Printf("Error 500")
			} else {
				fmt.Println("Unknown error occurred while fetching exchange rate.")
			}
		}

		salaryInUSD := person.GetSalaryAsFloat() * float64(exRate)
		if salaryInUSD > threshold {
			filtered = append(filtered, person)
		}
	}
	return filtered
}
