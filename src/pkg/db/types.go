package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type IntegerList []int64

func (a *IntegerList) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Int64Array: %v", value)
	}
	return json.Unmarshal(bytes, a)
}

func (a IntegerList) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

type StringList []string

func (a *StringList) Scan(value any) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringArray: %v", value)
	}
	return json.Unmarshal(bytes, a)
}

func (a StringList) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}
