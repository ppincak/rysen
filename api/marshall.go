package api

import (
	"encoding/json"
)

func Unmarshall(body []byte) (ApiResponse, error) {
	if body == nil {
		return nil, NewError("Body is null")
	}

	var m ApiResponse
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, NewError("Failed to unmarshall response")
	}
	return m, nil
}

func UnmarshallAs(body []byte, response interface{}) error {
	if body == nil {
		return NewError("Body is null")
	}
	err := json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	return nil
}
