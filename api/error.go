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
	if len(args) > 0 {
		message = fmt.Sprintf(message, args)
	}
	return Error{
		Message: message,
	}
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) toJsonString() string {
	b, _ := json.Marshal(e)
	return string(b)
}
