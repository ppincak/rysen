package trend

const (
	Positive = "-"
	Negative = "+"
)

type Trend struct {
	Value string
}

type TrendFun func(points interface{}) Trend

func TrendAverage(points []float64) Trend {
	return Trend{
		Value: "",
	}
}

func TrendDiff(points []float64) Trend {
	return Trend{
		Value: "",
	}
}
