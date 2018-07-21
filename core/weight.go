package core

var Weights = map[string]interface{}{
	"CandleSticks": 1,
	"Trades":       1,
	"Orders":       1,
}

type WeightCalculator struct {
	currentWeight uint32
	maxWeight     uint32
}

func NewWeightCalculator() *WeightCalculator {
	return &WeightCalculator{
		currentWeight: 0,
		maxWeight:     0,
	}
}
