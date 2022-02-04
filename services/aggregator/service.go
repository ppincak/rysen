package aggregator

import (
	"rysen/pkg/aggregate"
	b "rysen/pkg/bus"
	"rysen/pkg/errors"
	"rysen/pkg/scrape"
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
func (service *Service) Create(model *Model, aggregations aggregate.AggregationsMap) (*aggregate.Aggregator, error) {
	var conditionFunc aggregate.AggregationCondition
	convertedFunc, ok := aggregate.ConversionTable[model.Condition.Func]
	if !ok {
		return nil, errors.NewError("Aggregation function [%s] not found", model.Condition.Func)
	}

	switch convertedFunc {
	case aggregate.TillSize:
		conditionFunc = aggregate.AggretateTillSize(int(model.Condition.Value))
	case aggregate.TillTime:
		conditionFunc = aggregate.AggretateTillTime(model.Condition.Value)
	}

	aggregationFunc, ok := aggregations[aggregate.AggregationType(model.AggregationFunc)]
	if !ok {
		return nil, errors.NewError("Aggregation function not found")
	}

	aggregator := aggregate.NewAggregator(
		model.ReadTopic,
		model.WriteTopic,
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
