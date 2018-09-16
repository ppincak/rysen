package feed

type Metadata struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Topic             string `json:"topic"`
	EstimatedInterval int64  `json:"estimatedInterval"`
}

// Create Feed metadata
func NewMetadata(topic string, name string, description string) *Metadata {
	return &Metadata{
		Name:        name,
		Description: description,
		Topic:       topic,
	}
}
