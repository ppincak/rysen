package bus

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const BusBufferSize = 50

type uuidMap map[string]*busEventContainer
type topicMap map[string]uuidMap

type busEventContainer struct {
	uuid string
	outc chan *BusEvent
}

type Bus struct {
	topics topicMap
	uuids  map[string]string
	eventc chan *BusEvent
	subc   chan *BusSuscriptionEvent
	unsubc chan string
	stopc  chan struct{}
}

func NewBus() *Bus {
	return &Bus{
		topics: make(topicMap),
		uuids:  make(map[string]string),
		eventc: make(chan *BusEvent, BusBufferSize),
		subc:   make(chan *BusSuscriptionEvent),
		unsubc: make(chan string),
		stopc:  make(chan struct{}),
	}
}

func (bus *Bus) clean() {
	close(bus.stopc)
	close(bus.eventc)
}

func (bus *Bus) subscribe(subEvent *BusSuscriptionEvent) {
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

func (bus *Bus) unsubscribe(uuid string) {
	if topic, ok := bus.uuids[uuid]; ok {
		delete(bus.uuids, uuid)
		delete(bus.topics, topic)
	}
}

func (bus *Bus) Subscribe(topic string, outc chan *BusEvent) *BusSubscription {
	uuid := uuid.New().String()

	bus.subc <- &BusSuscriptionEvent{
		Topic: topic,
		Uuid:  uuid,
		Outc:  outc,
	}

	return &BusSubscription{
		Uuid: uuid,
		Cancel: func() {
			bus.unsubc <- uuid
		},
	}
}

func (bus *Bus) Unsubscribe(uuid string) {
	bus.unsubc <- uuid
}

func (bus *Bus) Publish(topic string, message interface{}) {
	bus.eventc <- &BusEvent{
		Topic:   topic,
		Message: message,
	}
}

func (bus *Bus) Start() {
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
			log.Infof("Stopping message bus")
			return
		}
	}
}

func (bus *Bus) Stop() {
	bus.stopc <- struct{}{}
}
