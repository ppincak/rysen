package services

import (
	"github.com/ppincak/rysen/api"
	b "github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/crypto"
	"github.com/ppincak/rysen/pkg/aggregate"
	"github.com/ppincak/rysen/pkg/scrape"
)

type AggregatorService struct {
	bus *b.Bus
}

func NewAggregatorService(bus *b.Bus) *AggregatorService {
	return &AggregatorService{
		bus: bus,
	}
}

func (service *AggregatorService) Create(metadata *AggregatorMetadata, exchange crypto.Exchange) (*aggregate.Aggregator, error) {
	var conditionFunc aggregate.AggregationCondition
	switch metadata.Condition.Func {
	case aggregate.TillSize:
		conditionFunc = aggregate.AggretateTillSize(int(metadata.Condition.Value))
	case aggregate.TillTime:
		conditionFunc = aggregate.AggretateTillTime(metadata.Condition.Value)
	}

	aggregations := exchange.Aggregations()
	aggregationFunc, ok := aggregations[aggregate.AggregationType(metadata.AggregationFunc)]
	if !ok {
		return nil, api.NewError("Aggregation function ot found")
	}

	aggregator := aggregate.NewAggregator(
		metadata.ReadTopic,
		metadata.WriteTopic,
		service.bus,
		ProcessCallerEvent,
		aggregationFunc,
		conditionFunc)

	go aggregator.Start()

	return aggregator, nil
}

// Caller event processing
func ProcessCallerEvent(event *b.BusEvent) (interface{}, error) {
	if assertion, ok := event.Message.(*scrape.CallerEvent); ok {
		return assertion.Data, nil
	} else {
		return nil, api.NewError("Failed to assert")
	}
}
