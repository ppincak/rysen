package bus

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ppincak/rysen/pkg/collections"
)

type AggregationResult struct {
	Result interface{} `json:"result"`
	From   int64       `json:"from"`
	To     int64       `json:"to"`
}

type Aggregator struct {
	bus                  *Bus
	sub                  *BusSubscription
	readTopic            string
	writeTopic           string
	eventc               chan *BusEvent
	stopc                chan struct{}
	list                 *collections.SliceList
	processFunc          ProcessFunc
	aggregationFunc      AggregationFunc
	aggregationCondition AggregationCondition
	from                 int64
	lastEntry            interface{}
}

type ProcessFunc func(event *BusEvent) (interface{}, error)
type AggregationCondition func(from int64, to int64, entries *collections.SliceList) bool

// Create new aggregator
func NewAggregator(
	readTopic string,
	writeTopic string,
	bus *Bus,
	processFunc ProcessFunc,
	aggregationFunc AggregationFunc,
	aggregationCondition AggregationCondition) *Aggregator {

	return &Aggregator{
		bus:                  bus,
		readTopic:            readTopic,
		writeTopic:           writeTopic,
		stopc:                make(chan struct{}),
		list:                 collections.NewSliceList(),
		processFunc:          processFunc,
		aggregationFunc:      aggregationFunc,
		aggregationCondition: aggregationCondition,
	}
}

// Start the aggregator, should be started as a goroutine
func (aggregator *Aggregator) Start() {
	defer func() {
		aggregator.sub.Cancel()

		log.Infof("Aggregator stopped  for topic [%s]", aggregator.readTopic)
	}()

	eventc := make(chan *BusEvent)
	aggregator.eventc = eventc
	aggregator.sub = aggregator.bus.Subscribe(aggregator.readTopic, eventc)

	log.Infof("Stated aggregator for topic [%s]", aggregator.readTopic)

	for {
		select {
		case event := <-aggregator.eventc:
			if aggregator.list.IsEmpty() {
				aggregator.from = time.Now().Unix()
			}
			to := time.Now().Unix()

			var lastEntry interface{}
			if processed, err := aggregator.processFunc(event); err == nil {
				aggregator.list.Add(processed)
				lastEntry = processed
			} else {
				continue
			}

			if aggregator.aggregationCondition(aggregator.from, to, aggregator.list) {
				// Note: might wan't to run this function in a separate goroutine
				result, err := aggregator.aggregationFunc(aggregator.list.EntriesCopy(), aggregator.lastEntry, aggregator.from)
				if err == nil {
					aggregator.bus.Publish(aggregator.writeTopic, &AggregationResult{
						Result: result,
						From:   aggregator.from,
						To:     to,
					})

					// last entry of this aggregation
					aggregator.lastEntry = lastEntry
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
