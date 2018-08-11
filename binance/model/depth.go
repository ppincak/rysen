package model

type OrdersDepth struct {
	Asks []interface{} `json:"asks"`
	Bids []interface{} `json:"bids"`
}
