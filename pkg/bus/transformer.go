package bus

import (
	log "github.com/sirupsen/logrus"
	"rysen/core"
)

// TODO should be renamed to subscriber !!

type BusEventHandler func(event *BusEvent)

type Transformer struct {
	eventc  chan *BusEvent
	outc    chan interface{}
	stopc   chan struct{}
	handler BusEventHandler
}

var _ core.Worker = (*Transformer)(nil)

func NewTransformer(eventc chan *BusEvent, handler BusEventHandler) *Transformer {
	return &Transformer{
		eventc:  eventc,
		stopc:   make(chan struct{}),
		handler: handler,
	}
}

func (transformer *Transformer) handle(event *BusEvent) {
	transformer.handler(event)
}

func (transformer *Transformer) Start() {
	defer func() {
		log.Info("Transformer worker finished")
	}()

	for {
		select {
		case event := <-transformer.eventc:
			transformer.handle(event)
		case <-transformer.stopc:
			return
		}
	}
}

func (transformer *Transformer) Stop() {
	transformer.stopc <- struct{}{}
}
