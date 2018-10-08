package publisher

import (
	"context"

	b "github.com/ppincak/rysen/pkg/bus"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	*Model

	brokers      []string
	writer       *kafka.Writer
	subscription *b.BusSubscription
	transformer  *b.Transformer
}

func NewKafkaPublisher(model *Model, brokers []string) *KafkaPublisher {
	return &KafkaPublisher{
		Model:   model,
		brokers: brokers,
	}
}

// Initialize the publisher
func (publisher *KafkaPublisher) Init(bus *b.Bus) {
	publisher.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  publisher.brokers,
		Topic:    publisher.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})

	outc := make(chan *b.BusEvent)
	publisher.subscription = bus.Subscribe(publisher.ReadTopic, outc)
	publisher.transformer = b.NewTransformer(outc, publisher.publishHandler)

	go publisher.transformer.Start()
}

// Destroy publisher
func (publisher *KafkaPublisher) Destroy() {
	publisher.transformer.Stop()
}

// Return topic
func (publisher *KafkaPublisher) Topic() string {
	return publisher.KafkaTopic
}

// Publish to kafka broker
func (publisher *KafkaPublisher) Publish(message interface{}) {
	publisher.writer.WriteMessages(context.Background())
}

// Handler for publishing
func (publisher *KafkaPublisher) publishHandler(event *b.BusEvent) {
	publisher.Publish(event)
}
