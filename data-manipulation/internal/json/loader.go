package json

import (
	"data-manipulation/internal/models"
	"encoding/json"
	"io/ioutil"
)

// LoadPersons reads a JSON file and unmarshals its content into a Persons struct.
func LoadPersons(filePath string) (*models.Persons, error) {
	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal into the Persons struct
	var persons models.Persons
	err = json.Unmarshal(data, &persons)
	if err != nil {
		return nil, err
	}

	return &persons, nil
}
