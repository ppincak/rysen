package services

type FeedMetadata struct {
	Name        string
	Description string
	Topic       string
}

func NewFeedMetadata(topic string, name string, description string) *FeedMetadata {
	return &FeedMetadata{
		Name:        name,
		Description: description,
		Topic:       topic,
	}
}
