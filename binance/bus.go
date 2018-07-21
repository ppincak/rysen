package binance

import (
	"github.com/google/uuid"
	"github.com/ppincak/rysen/core"
	log "github.com/sirupsen/logrus"
)

// TODO move to common ?
const BinanceBusBufferSize = 50

type uuidMap map[string]*busEventContainer
type topicMap map[string]uuidMap

type busEventContainer struct {
	uuid string
	outc chan *core.BusEvent
}

type BinanceBus struct {
	topics topicMap
	uuids  map[string]string
	eventc chan *core.BusEvent
	subc   chan *core.BusSuscriptionEvent
	unsubc chan string
	stopc  chan struct{}
}

func NewBinanceBus() *BinanceBus {
	return &BinanceBus{
		topics: make(topicMap),
		uuids:  make(map[string]string),
		eventc: make(chan *core.BusEvent, BinanceBusBufferSize),
		subc:   make(chan *core.BusSuscriptionEvent),
		unsubc: make(chan string),
		stopc:  make(chan struct{}),
	}
}

func (bus *BinanceBus) clean() {
	close(bus.stopc)
	close(bus.eventc)
}

func (bus *BinanceBus) subscribe(subEvent *core.BusSuscriptionEvent) {
	var subs uuidMap
	if val, ok := bus.topics[subEvent.Topic]; ok {
		subs = val
	} else {
		subs = make(uuidMap)
	}

	event := &busEventContainer{
		uuid: subEvent.Uuid,
		outc: subEvent.Outc,
	}

	subs[subEvent.Uuid] = event
	bus.uuids[subEvent.Uuid] = subEvent.Topic
	bus.topics[subEvent.Topic] = subs
}

func (bus *BinanceBus) unsubscribe(uuid string) {
	if topic, ok := bus.uuids[uuid]; ok {
		delete(bus.uuids, uuid)
		delete(bus.topics, topic)
	}
}

func (bus *BinanceBus) Subscribe(topic string, outc chan *core.BusEvent) *core.BusSubscription {
	uuid := uuid.New().String()

	bus.subc <- &core.BusSuscriptionEvent{
		Topic: topic,
		Uuid:  uuid,
		Outc:  outc,
	}

	return &core.BusSubscription{
		Uuid: uuid,
		Cancel: func() {
			bus.unsubc <- uuid
		},
	}
}

func (bus *BinanceBus) Unsubscribe(uuid string) {
	bus.unsubc <- uuid
}

func (bus *BinanceBus) Publish(topic string, message interface{}) {
	bus.eventc <- &core.BusEvent{
		Topic:   topic,
		Message: message,
	}
}

func (bus *BinanceBus) Start() {
	defer bus.clean()

	for {
		select {
		case event := <-bus.eventc:
			if subsribers, ok := bus.topics[event.Topic]; ok {
				for _, subsriber := range subsribers {
					subsriber.outc <- event
				}
			}
		case event := <-bus.subc:
			bus.subscribe(event)
		case uuid := <-bus.unsubc:
			bus.unsubscribe(uuid)
		case <-bus.stopc:
			log.Infof("Stopping Binance message bus")
			return
		}
	}
}

func (bus *BinanceBus) Stop() {
	bus.stopc <- struct{}{}
}
