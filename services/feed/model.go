package feed

type Model struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Topic             string `json:"topic"`
	EstimatedInterval int64  `json:"estimatedInterval"`
}

// Create Feed model
func NewModel(topic string, name string, description string) *Model {
	return &Model{
		Name:        name,
		Description: description,
		Topic:       topic,
	}
}
