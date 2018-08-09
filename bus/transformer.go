package bus

import (
	"github.com/ppincak/rysen/core"
	log "github.com/sirupsen/logrus"
)

// TODO: Rename to interceptor func
type TransformFuc func(event *BusEvent)

type Transformer struct {
	eventc        chan *BusEvent
	outc          chan interface{}
	stopc         chan struct{}
	transformFunc TransformFuc
}

var _ core.Worker = (*Transformer)(nil)

func NewTransformer(eventc chan *BusEvent, transformFunc TransformFuc) *Transformer {
	return &Transformer{
		eventc:        eventc,
		stopc:         make(chan struct{}),
		transformFunc: transformFunc,
	}
}

func (transformer *Transformer) transform(event *BusEvent) {
	transformer.transformFunc(event)
}

func (transformer *Transformer) Start() {
	defer func() {
		log.Info("Transformer worker finished")
	}()

	for {
		select {
		case event := <-transformer.eventc:
			transformer.transform(event)
		case <-transformer.stopc:
			return
		}
	}
}

func (transformer *Transformer) Stop() {
	transformer.stopc <- struct{}{}
}
