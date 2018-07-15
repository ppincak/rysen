package common

type Statistic struct {
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
}

func NewStatistic(name string) *Statistic {
	return &Statistic{
		Name:   name,
		Values: make(map[string]interface{}),
	}
}
