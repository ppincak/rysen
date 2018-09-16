package math

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/floats"

	"gonum.org/v1/gonum/stat"
)

// Contains different averages
type Average struct {
	Mean     float64 `json:"mean"`
	Median   float64 `json:"median"`
	Variance float64 `json:"variance"`
	Stddev   float64 `json:"stddev"`
	Sum      float64 `json:"sum"`
	Length   int     `json:"length"`
}

// Create new Average
func NewAverage(values []float64) *Average {
	sort.Float64s(values)

	variance := stat.Variance(values, nil)

	return &Average{
		Mean:     stat.Mean(values, nil),
		Variance: variance,
		Median:   stat.Quantile(0.5, stat.Empirical, values, nil),
		Stddev:   math.Sqrt(variance),
		Sum:      floats.Sum(values),
		Length:   len(values),
	}
}
