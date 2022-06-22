package app

import (
	"encoding/json"
	"fmt"
	"homework-l0/internal/models"
)

//Unmarshal данных считаных с канала nats-streaming, отсеивание мусора
func ParseMessages(data []byte) (*models.Orders, error) {
	ord := models.Orders{}
	err := json.Unmarshal(data, &ord)
	if err != nil {
		return nil, err
	}
	if ord.Entry != "WBIL" {
		return nil, fmt.Errorf("wrong or missing entry")
	}
	return &ord, nil
}
