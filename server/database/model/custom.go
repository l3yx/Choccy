package model

import (
	"database/sql/driver"
	"encoding/json"
)

type StrArr []string

func (s *StrArr) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, s)
}

func (s StrArr) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type StrMap map[string]interface{}

func (s *StrMap) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, s)
}

func (s StrMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type IsRead struct {
	IdList []int `json:"idList"`
	Read   bool  `json:"read"`
}
