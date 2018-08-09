package data

import (
	"strconv"

	"github.com/ppincak/rysen/api"
)

func ExtractPrices(list []interface{}) ([]float64, error) {
	result := make([]float64, len(list))
	for i, value := range list {
		if asserted, ok := value.(api.ApiResponse); ok {
			if price, ok := asserted["price"]; ok {
				switch t := price.(type) {
				case float32:
					result[i] = float64(t)
				case float64:
					result[i] = t
				case string:
					result[i], _ = strconv.ParseFloat(t, 64)
				}
			} else {
				return nil, api.NewError("Map doesn't contain price")
			}
		} else {
			return nil, api.NewError("Type assertion failed")
		}
	}
	return result, nil
}
