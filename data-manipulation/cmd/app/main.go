package main

import (
	"context"
	"data-manipulation/internal/appbase"
	"data-manipulation/internal/clients/cadana_backend"
	cadanaAPI "data-manipulation/internal/clients/cadana_backend/openapi"
	"data-manipulation/internal/json"
	"data-manipulation/internal/services"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

const (
	serviceName = "data-manipulation.server"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}

	ctx, mainCtxStop := context.WithCancel(context.Background())

	application := appbase.New(
		appbase.Init(serviceName),
		appbase.WithSignals(ctx, mainCtxStop),
		appbase.WithSentry(),
	)

	// Path to JSON file
	filePath := "data/persons.json"

	// Load data from JSON file
	personsData, err := json.LoadPersons(filePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load persons data")
	}

	var (
		cadanaHandler = initCadanaService(cfg, application)
	)

	// Initialize the ServicePersons service with data
	personsService := services.ServicePersons{Data: personsData.Data}

	// Demonstrate functionalities
	fmt.Println("Original Data:")
	fmt.Println(personsService.Data)

	// Sort by salary ascending
	personsService.SortBySalaryAscending()
	fmt.Println("\nSorted by Salary (Ascending):")
	fmt.Println(personsService.Data)

	// Sort by salary descending
	personsService.SortBySalaryDescending()
	fmt.Println("\nSorted by Salary (Descending):")
	fmt.Println(personsService.Data)

	// Group by currency
	grouped := personsService.GroupByCurrency()
	fmt.Println("\nGrouped by Currency:")
	for currency, group := range grouped {
		fmt.Printf("%s: %+v\n", currency, group)
	}

	// Filter by salary in USD (threshold 15 USD)
	filtered := personsService.FilterBySalaryInUSD(15, *cadanaHandler)
	fmt.Println("\nFiltered by Salary > 15 USD:")
	fmt.Println(filtered)
}

func initCadanaService(cfg *Config, application *appbase.AppBase) *cadana_backend.APIClient {
	cadanaClient, err := cadanaAPI.NewClientWithResponses(
		cfg.CadanaBackendBaseURL,
	)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal().Err(err).Msg("failed to initialize gotron service client")
	}

	return cadana_backend.NewAPIClient(cadanaClient)
}
