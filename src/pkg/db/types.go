package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
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

type IPList []net.IP

func (a *IPList) Scan(value any) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan IPArray: %v", value)
	}

	var strArr []string
	if err := json.Unmarshal(bytes, &strArr); err != nil {
		return fmt.Errorf("failed to unmarshal IPArray: %w", err)
	}

	ips := make([]net.IP, 0, len(strArr))
	for _, s := range strArr {
		ip := net.ParseIP(s)
		if ip == nil {
			return fmt.Errorf("invalid IP address in JSON: %s", s)
		}
		ips = append(ips, ip)
	}
	*a = ips
	return nil
}

func (a IPList) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	strArr := make([]string, 0, len(a))
	for _, ip := range a {
		strArr = append(strArr, ip.String())
	}
	return json.Marshal(strArr)
}
