package bus

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ppincak/rysen/pkg/collections"
)

type AggregationResult struct {
	Result interface{}
	From   int64
	To     int64
}

type Aggregator struct {
	bus                  *Bus
	sub                  *BusSubscription
	topic                string
	resultTopic          string
	eventc               chan *BusEvent
	stopc                chan struct{}
	list                 *collections.SliceList
	processFunc          ProcessFunc
	aggregationFunc      AggregationFunc
	aggregationCondition AggregationCondition
	from                 int64
}

type ProcessFunc func(event *BusEvent) (interface{}, error)
type AggregationCondition func(from int64, to int64, size int) bool

// Create new aggregator
func NewAggregator(
	topic string,
	suffix string,
	bus *Bus,
	processFunc ProcessFunc,
	aggregationFunc AggregationFunc,
	aggregationCondition AggregationCondition) *Aggregator {

	return &Aggregator{
		bus:                  bus,
		topic:                topic,
		resultTopic:          strings.Join([]string{topic, suffix}, "-"),
		stopc:                make(chan struct{}),
		list:                 collections.NewSliceList(10),
		processFunc:          processFunc,
		aggregationFunc:      aggregationFunc,
		aggregationCondition: aggregationCondition,
	}
}

// Start the aggregator, should be started as a goroutine
func (aggregator *Aggregator) Start() {
	defer func() {
		aggregator.sub.Cancel()

		log.Infof("Aggregator stopped  for topic [%s]", aggregator.topic)
	}()

	eventc := make(chan *BusEvent)
	aggregator.eventc = eventc
	aggregator.sub = aggregator.bus.Subscribe(aggregator.topic, eventc)

	log.Infof("Stated aggregator for topic [%s]", aggregator.topic)

	for {
		select {
		case event := <-aggregator.eventc:
			if aggregator.list.IsEmpty() {
				aggregator.from = time.Now().Unix()
			}
			to := time.Now().Unix()

			if processed, err := aggregator.processFunc(event); err == nil {
				aggregator.list.Add(processed)
			} else {
				continue
			}

			if aggregator.aggregationCondition(aggregator.from, to, aggregator.list.Size()) {
				result, err := aggregator.aggregationFunc(aggregator.list.EntriesCopy())
				if err == nil {
					aggregator.bus.Publish(aggregator.resultTopic, &AggregationResult{
						Result: result,
						From:   aggregator.from,
						To:     to,
					})

					aggregator.list.Reset()
				} else {
					log.Error("Error during aggregation")
					log.Error(err)
				}
			}
		case <-aggregator.stopc:
			return
		}
	}
}

// Stop the aggregator, should be started as a goroutine
func (aggregator *Aggregator) Stop() {
	aggregator.stopc <- struct{}{}
}
