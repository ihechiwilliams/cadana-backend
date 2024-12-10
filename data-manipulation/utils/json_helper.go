package utils

import (
	"data-manipulation/models"
	"encoding/json"
	"io/ioutil"
)

func LoadPersons(filePath string) (*models.Persons, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var persons models.Persons
	err = json.Unmarshal(data, &persons)
	if err != nil {
		return nil, err
	}
	return &persons, nil
}
