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

// func ExtractOrders(list []interface{}) (interface{}, error) {
// 	result := make([]float64, len(list))
// 	for i, value := range list {
// 		if response, ok := value.(api.ApiResponse); ok {
// 			extractOrder("asks", response)
// 			extractOrder("bids", response)
// 		} else {
// 			return nil, api.NewError("Type assertion failed")
// 		}
// 	}
// 	return result, nil
// }

// func extractOrder(property string, response api.ApiResponse) {
// 	if value, ok := response[property]; ok {
// 		if list, ok := value.([]interface{}); ok {
// 			for i, value := range list {
// 				if order, ok := value.([]string); ok {
// 					price, _ := toFloat64(order[0])
// 					volume, _ := toFloat64(order[1])
// 				}
// 			}
// 		}
// 	}
// }

func toFloat64(value interface{}) (float64, error) {
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

	return 0, api.NewError("")
}
