package services

import (
	b "github.com/ppincak/rysen/bus"
)

type AggregatorService struct {
	bus *b.Bus
}

func NewAggregatorService(bus *b.Bus) *AggregatorService {
	return &AggregatorService{
		bus: bus,
	}
}

func (service *AggregatorService) Create(metadata *AggregatorMetadata) *b.Aggregator {
	aggregator := b.NewAggregator(metadata.ReadTopic, metadata.WriteTopic, service.bus, nil, nil, nil)
	go aggregator.Start()
	return aggregator
}
