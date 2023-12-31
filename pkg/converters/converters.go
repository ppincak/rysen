package converters

import (
	"strconv"

	"rysen/pkg/errors"
)

// Convert value to float64
func ToFloat64(value interface{}) (float64, error) {
	switch t := value.(type) {
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	case string:
		parsed, err := strconv.ParseFloat(t, 64)
		if err == nil {
			return parsed, nil
		}
		return 0, err
	}

	return 0, errors.NewError("Failed to convert")
}
