package math

import (
	"gonum.org/v1/gonum/floats"
)

type Sum struct {
	Metric float64 `json:"metric"`
	Volume float64 `json:"volume"`
}

// New sum
func NewSum(metric float64, volume float64) *Sum {
	return &Sum{
		Metric: metric,
		Volume: volume,
	}
}

// Summed values grouped by float64 key
func SumGroupedByFloat64(entries map[float64][]float64) []*Sum {
	sums := make([]*Sum, len(entries))
	i := 0
	for key, values := range entries {
		sums[i] = &Sum{
			Metric: key,
			Volume: floats.Sum(values),
		}
		i++
	}
	return sums
}
