package app

import (
	"encoding/json"
	"homework-l0/internal/models"
)

func ParseMessages(data []byte) (*models.Orders, error) {
	ord := models.Orders{}
	err := json.Unmarshal(data, &ord)
	if err != nil {
		return nil, err
	}
	if ord.Entry != "WBIL" {
		return nil, err
	}

	return &ord, nil
}
