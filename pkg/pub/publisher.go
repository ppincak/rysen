package pub

type Publisher interface {
	Topic() string
	Publish(message interface{})
}
