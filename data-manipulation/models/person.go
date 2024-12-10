package models

import (
	"sort"
	"strconv"
)

type Person struct {
	ID         string `json:"id"`
	PersonName string `json:"personName"`
	Salary     Salary `json:"salary"`
}

type Salary struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

func (p *Person) GetSalaryAsFloat() float64 {
	val, _ := strconv.ParseFloat(p.Salary.Value, 64)
	return val
}

type Persons struct {
	Data []Person `json:"data"`
}

func (ps *Persons) SortBySalaryAscending() {
	sort.Slice(ps.Data, func(i, j int) bool {
		return ps.Data[i].GetSalaryAsFloat() < ps.Data[j].GetSalaryAsFloat()
	})
}

func (ps *Persons) SortBySalaryDescending() {
	sort.Slice(ps.Data, func(i, j int) bool {
		return ps.Data[i].GetSalaryAsFloat() > ps.Data[j].GetSalaryAsFloat()
	})
}

func (ps *Persons) GroupByCurrency() map[string][]Person {
	grouped := make(map[string][]Person)
	for _, person := range ps.Data {
		currency := person.Salary.Currency
		grouped[currency] = append(grouped[currency], person)
	}
	return grouped
}

func (ps *Persons) FilterBySalaryInUSD(threshold float64, exchangeRates ExchangeRateProvider) []Person {
	var filtered []Person
	for _, person := range ps.Data {
		rate := exchangeRates.GetRate(person.Salary.Currency, "USD")
		salaryInUSD := person.GetSalaryAsFloat() * rate
		if salaryInUSD > threshold {
			filtered = append(filtered, person)
		}
	}
	return filtered
}

type ExchangeRateProvider interface {
	GetRate(from, to string) float64
}
