package main

import (
	"data-manipulation/utils"
	"fmt"
)

func main() {
	// Load JSON data
	filePath := "data/persons.json"
	persons, err := utils.LoadPersons(filePath)
	if err != nil {
		fmt.Printf("Error loading persons: %v\n", err)
		return
	}

	// Display original data
	fmt.Println("Original Data:", persons.Data)

	// Sort by salary ascending
	persons.SortBySalaryAscending()
	fmt.Println("Sorted Ascending:", persons.Data)

	// Sort by salary descending
	persons.SortBySalaryDescending()
	fmt.Println("Sorted Descending:", persons.Data)

	// Group by currency
	grouped := persons.GroupByCurrency()
	fmt.Println("Grouped by Currency:", grouped)

	// Filter by salary in USD
	exchangeRates := utils.NewExchangeRates("https://api.example.com/v1/exchange-rate") // Simulated rates
	filtered := persons.FilterBySalaryInUSD(15, exchangeRates)
	fmt.Println("Filtered by Salary in USD > 15:", filtered)
}
