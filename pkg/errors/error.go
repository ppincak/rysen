package errors

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// create ne error
func NewError(message string, args ...interface{}) Error {
	err := Error{}
	err.Format(message, args)
	return err
}

// create ne error with error code
func NewErrorWithCode(message string, code string) Error {
	return Error{
		Message: message,
		Code:    code,
	}
}

// get error mesage
func (e Error) Format(message string, args ...interface{}) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args)
	}
	e.Message = message
}

// get error mesage
func (e Error) Error() string {
	return e.Message
}

func (e Error) toJsonString() string {
	b, _ := json.Marshal(e)
	return string(b)
}
