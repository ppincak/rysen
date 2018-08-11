package services

type FeedMetadata struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Topic             string `json:"topic"`
	EstimatedInterval int64  `json:"estimatedInterval"`
}

// Create feed metadata
func NewFeedMetadata(topic string, name string, description string) *FeedMetadata {
	return &FeedMetadata{
		Name:        name,
		Description: description,
		Topic:       topic,
	}
}
