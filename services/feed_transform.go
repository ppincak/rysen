package services

import (
	"reflect"

	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/pkg/scrape"
)

func TransformForWsClient(message interface{}) (interface{}, error) {
	switch message.(type) {
	case *scrape.CallerEvent:
		return message.(*scrape.CallerEvent).Data, nil
	case *bus.AggregationResult:
		return message.(*bus.AggregationResult).Result, nil
	}
	return nil, api.NewError("Unhandled type [%s]", reflect.TypeOf(message).Name)
}
