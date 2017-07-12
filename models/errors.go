package models

import (
	"github.com/asaskevich/govalidator"
)

type Errors interface {
	ErrorsAsMap() map[string]string
}

type GovalidatorErrors struct {
	Err error
}

func (g GovalidatorErrors) ErrorsAsMap() map[string]string {
	return govalidator.ErrorsByField(g.Err)
}
