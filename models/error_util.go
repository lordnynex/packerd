package models

import (
	"fmt"
)

func (m *Error) String() string {
	return fmt.Sprintf("Error: [%d %s]", m.Code, m.Message)
}

func NewError(code int, message string) *Error {
	e := new(Error)
	e.Code = int64(code)
	e.Message = &message
	return e
}
