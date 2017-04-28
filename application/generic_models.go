package application

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

/**
 *  remember to use `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
 */

type uintList []uint

func (u uintList) Value() (driver.Value, error) {
	return value(u)
}

func (u *uintList) Scan(src interface{}) error {
	return scan(u, src)
}

type stringList []string

func (s stringList) Value() (driver.Value, error) {
	return value(s)
}

func (s *stringList) Scan(src interface{}) error {
	return scan(s, src)
}

func value(entity interface{}) (driver.Value, error) {
	val, err := json.Marshal(entity)
	return val, err
}

func scan(entity, src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	err := json.Unmarshal(source, entity)
	if err != nil {
		return err
	}

	return nil
}

func pValue(entity interface{}) (driver.Value, error) {
	var josnb JSONB
	josnb, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	return josnb.Value()
}

func pScan(entity, src interface{}) error {
	var jsonb JSONB
	err := (&jsonb).Scan(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonb, entity)
}

type JSONB []byte

func (j JSONB) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		errors.New("Scan source was not string")
	}
	// I think I need to make a copy of the bytes.
	// It seems the byte slice passed in is re-used
	*j = append((*j)[0:0], s...)

	return nil
}

// MarshalJSON returns *m as the JSON encoding of m.
func (m JSONB) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *JSONB) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (j JSONB) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSONB) Equals(j1 JSONB) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}
