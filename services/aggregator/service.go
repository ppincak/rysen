package aggregator

import (
	"github.com/ppincak/rysen/pkg/aggregate"
	b "github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/errors"
	"github.com/ppincak/rysen/pkg/scrape"
)

type Service struct {
	bus *b.Bus
}

// Create new Servie
func NewService(bus *b.Bus) *Service {
	return &Service{
		bus: bus,
	}
}

// Note: think of better implementation
// Create aggregator
func (service *Service) Create(metadata *Metadata, aggregations aggregate.AggregationsMap) (*aggregate.Aggregator, error) {
	var conditionFunc aggregate.AggregationCondition
	switch metadata.Condition.Func {
	case aggregate.TillSize:
		conditionFunc = aggregate.AggretateTillSize(int(metadata.Condition.Value))
	case aggregate.TillTime:
		conditionFunc = aggregate.AggretateTillTime(metadata.Condition.Value)
	}

	aggregationFunc, ok := aggregations[aggregate.AggregationType(metadata.AggregationFunc)]
	if !ok {
		return nil, errors.NewError("Aggregation function ot found")
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
		return nil, errors.NewError("Failed to assert")
	}
}
