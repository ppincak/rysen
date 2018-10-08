package feed

import (
	"reflect"

	"github.com/ppincak/rysen/pkg/aggregate"
	"github.com/ppincak/rysen/pkg/errors"
	"github.com/ppincak/rysen/pkg/scrape"
)

// Do message trasformation
func TransformForWsClient(message interface{}) (interface{}, error) {
	switch message.(type) {
	case *scrape.CallerEvent:
		return message.(*scrape.CallerEvent).Data, nil
	case *aggregate.AggregationResult:
		return message.(*aggregate.AggregationResult).Result, nil
	}
	return nil, errors.NewError("Unhandled type [%s]", reflect.TypeOf(message).Name)
}
