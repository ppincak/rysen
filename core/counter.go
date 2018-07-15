package core

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const counterEventBufferSize = 50

const (
	resetRequests uint32 = iota
	resetOrders
	getRequests
	getOrders
	incrementRequests
	incrementOrders
)

// Note: add notification option
type ApiCallCounter struct {
	OrdersMade    int64
	RequestsMade  int64
	ordersTicker  *time.Ticker
	requestTicker *time.Ticker
	ordersLimit   int64
	requestsLimit int64
	evc           chan *CountEvent
	stopc         chan struct{}
}

type CountEvent struct {
	Event uint32
	resc  chan<- int64
}

func NewApiCallCounter(orderMillis int64, requestMillis int64) *ApiCallCounter {
	return &ApiCallCounter{
		OrdersMade:    0,
		RequestsMade:  0,
		ordersTicker:  time.NewTicker(ToDuration(orderMillis)),
		requestTicker: time.NewTicker(ToDuration(requestMillis)),
		evc:           make(chan *CountEvent, counterEventBufferSize),
		stopc:         make(chan struct{}),
	}
}

func (counter *ApiCallCounter) Start() {
	for {
		select {
		case <-counter.ordersTicker.C:
			log.Debugf("Orders counter has been reset. Requests made: %d", counter.OrdersMade)
			counter.OrdersMade = 0
		case <-counter.requestTicker.C:
			log.Debugf("Requests counter has been reset. Requests made: %d", counter.RequestsMade)
			counter.RequestsMade = 0
		case msg := <-counter.evc:
			switch msg.Event {
			case getOrders:
				msg.resc <- counter.OrdersMade
			case getRequests:
				msg.resc <- counter.RequestsMade
			case resetOrders:
				counter.OrdersMade = 0
			case resetRequests:
				counter.RequestsMade = 0
			case incrementOrders:
				counter.OrdersMade++
			case incrementRequests:
				counter.RequestsMade++
			}
		case <-counter.stopc:
			return
		}
	}
}

func (counter *ApiCallCounter) Stop() {
	counter.stopc <- struct{}{}
}

func (counter *ApiCallCounter) sendAndGet(resc chan int64, event uint32) int64 {
	counter.evc <- &CountEvent{
		Event: event,
		resc:  resc,
	}
	if resc == nil {
		return -1
	}
	return <-resc
}

func (counter *ApiCallCounter) IncrRequests(resc chan int64) int64 {
	return counter.sendAndGet(resc, incrementRequests)
}

func (counter *ApiCallCounter) IncrOrders(resc chan int64) int64 {
	return counter.sendAndGet(resc, incrementOrders)
}

func (counter *ApiCallCounter) ResetRequests(resc chan int64) {
	counter.evc <- &CountEvent{
		Event: resetRequests,
	}
}

func (counter *ApiCallCounter) ResetOrders(resc chan int64) {
	counter.evc <- &CountEvent{
		Event: resetOrders,
	}
}
