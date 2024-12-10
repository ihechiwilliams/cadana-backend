package models

import "strconv"

type Persons struct {
	Data []Person `json:"data"`
}

type Person struct {
	ID         string `json:"id"`
	PersonName string `json:"personName"`
	Salary     Salary `json:"salary"`
}

func (p *Person) GetSalaryAsFloat() float64 {
	val, _ := strconv.ParseFloat(p.Salary.Value, 64)
	return val
}

type Salary struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}
