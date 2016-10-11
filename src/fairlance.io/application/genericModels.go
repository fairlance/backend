package application

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

/**
 *  remember to use `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
 */

type uints []uint

func (u uints) Value() (driver.Value, error) {
	return value(u)
}

func (u *uints) Scan(src interface{}) error {
	return scan(u, src)
}

type strings []string

func (s strings) Value() (driver.Value, error) {
	return value(s)
}

func (s *strings) Scan(src interface{}) error {
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
