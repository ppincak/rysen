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
type CallLimiter struct {
	OrdersMade    int64
	RequestsMade  int64
	ordersTicker  *time.Ticker
	requestTicker *time.Ticker
	ordersLimit   int64
	requestsLimit int64
	evc           chan *limiterEvent
	stopc         chan struct{}
}

type limiterEvent struct {
	event uint32
	resc  chan<- int64
}

func NewApiCallCounter(orderMillis int64, requestMillis int64) *CallLimiter {
	return &CallLimiter{
		OrdersMade:    0,
		RequestsMade:  0,
		ordersTicker:  time.NewTicker(ToDurationMillis(orderMillis)),
		requestTicker: time.NewTicker(ToDurationMillis(requestMillis)),
		evc:           make(chan *limiterEvent, counterEventBufferSize),
		stopc:         make(chan struct{}),
	}
}

func (limiter *CallLimiter) Start() {
	for {
		select {
		case <-limiter.ordersTicker.C:
			log.Debugf("Orders counter has been reset. Requests made: %d", limiter.OrdersMade)
			limiter.OrdersMade = 0
		case <-limiter.requestTicker.C:
			log.Debugf("Requests counter has been reset. Requests made: %d", limiter.RequestsMade)
			limiter.RequestsMade = 0
		case msg := <-limiter.evc:
			switch msg.event {
			case getOrders:
				msg.resc <- limiter.OrdersMade
			case getRequests:
				msg.resc <- limiter.RequestsMade
			case resetOrders:
				limiter.OrdersMade = 0
			case resetRequests:
				limiter.RequestsMade = 0
			case incrementOrders:
				limiter.OrdersMade++
			case incrementRequests:
				limiter.RequestsMade++
			}
		case <-limiter.stopc:
			return
		}
	}
}

func (limiter *CallLimiter) Stop() {
	limiter.stopc <- struct{}{}
}

func (limiter *CallLimiter) sendAndGet(resc chan int64, event uint32) int64 {
	limiter.evc <- &limiterEvent{
		event: event,
		resc:  resc,
	}
	if resc == nil {
		return -1
	}
	return <-resc
}

func (limiter *CallLimiter) IncrRequests(resc chan int64) int64 {
	return limiter.sendAndGet(resc, incrementRequests)
}

func (limiter *CallLimiter) IncrOrders(resc chan int64) int64 {
	return limiter.sendAndGet(resc, incrementOrders)
}

func (limiter *CallLimiter) ResetRequests(resc chan int64) {
	limiter.evc <- &limiterEvent{
		event: resetRequests,
	}
}

func (limiter *CallLimiter) ResetOrders(resc chan int64) {
	limiter.evc <- &limiterEvent{
		event: resetOrders,
	}
}
