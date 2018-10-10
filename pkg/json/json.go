package json

import (
	"encoding/json"

	"github.com/ppincak/rysen/pkg/errors"
)

// unmarshall as mao
func Unmarshall(body []byte) (map[string]interface{}, error) {
	if body == nil {
		return nil, errors.NewError("Body is null")
	}

	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, errors.NewError("Failed to unmarshall response")
	}
	return m, nil
}

// unmarshall as given type
func UnmarshallAs(body []byte, response interface{}) error {
	if body == nil {
		return errors.NewError("Body is null")
	}
	err := json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	return nil
}
