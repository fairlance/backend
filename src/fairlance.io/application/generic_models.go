package application

import (
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
