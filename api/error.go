package api

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func NewError(message string, args ...interface{}) Error {
	return Error{
		Message: fmt.Sprintf(message, args),
	}
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) toJsonString() string {
	b, _ := json.Marshal(e)
	return string(b)
}
