package core

type Bus interface {
	// Subscribe to topic with channel
	Subscribe(topic string, outc chan *BusEvent) *BusSubscription
	// Unsubscribe from topic with channel
	Unsubscribe(uuid string)
}

// Event which us sent through implementation of Bus interface
type BusEvent struct {
	Topic   string
	Message interface{}
}

// Event for use internally by event bus
type BusSuscriptionEvent struct {
	Topic string
	Uuid  string
	Outc  chan *BusEvent
}

// Subscription created after subscribing a topic
type BusSubscription struct {
	Uuid   string
	Cancel func()
}
