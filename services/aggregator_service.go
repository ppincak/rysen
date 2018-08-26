package services

import (
	b "github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/pkg/aggregate"
)

type AggregatorService struct {
	bus *b.Bus
}

func NewAggregatorService(bus *b.Bus) *AggregatorService {
	return &AggregatorService{
		bus: bus,
	}
}

func (service *AggregatorService) Create(metadata *AggregatorMetadata) *aggregate.Aggregator {
	// var conditionFunc aggregate.AggregationCondition
	// switch metadata.Condition.Func {
	// case aggregate.TillSize:
	// 	conditionFunc = aggregate.AggretateTillSize(int(metadata.Condition.Value))
	// case aggregate.TillTime:
	// 	conditionFunc = aggregate.AggretateTillTime(metadata.Condition.Value)
	// }

	aggregator := aggregate.NewAggregator(metadata.ReadTopic, metadata.WriteTopic, service.bus, nil, nil, nil)

	go aggregator.Start()
	return aggregator
}
