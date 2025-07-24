package models

import (
	"database/sql/driver"
	"fmt"

	json "github.com/goccy/go-json"
)

type Address struct {
	Type   string `json:"type"` // e.g., "home", "work"
	Street string `json:"street"`
	City   string `json:"city"`
}

type AddressList []Address

func (a AddressList) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AddressList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	return json.Unmarshal(b, a)
}
Write the seed code add the seed command to the goapi 

