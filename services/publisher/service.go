package publisher

import (
	"strings"
	"sync"

	"rysen/pkg/bus"
	"rysen/pkg/errors"
	"rysen/pkg/pub"
)

type Service struct {
	bus          *bus.Bus
	kafkaBrokers []string

	publishers map[string]*KafkaPublisher
	lock       *sync.RWMutex
}

// new service
func NewService(bus *bus.Bus, kafkaBrokers []string) *Service {
	return &Service{
		bus:          bus,
		kafkaBrokers: kafkaBrokers,
	}
}

// assemble key
func (service *Service) assembleKey(model *Model) string {
	return strings.Join([]string{model.ReadTopic, model.KafkaTopic}, "/")
}

// create publisher
func (service *Service) CreatePublisher(model *Model) (pub.Publisher, error) {
	defer service.lock.Unlock()
	service.lock.Unlock()

	key := service.assembleKey(model)

	if _, ok := service.publishers[key]; ok {
		return nil, errors.NewError("Publisher [%s] already exists", key)
	}
	publisher := NewKafkaPublisher(model, nil)
	service.publishers[key] = publisher

	return publisher, nil
}

// get list of publishers
func (service *Service) ListPublishers() []*Model {
	defer service.lock.RUnlock()
	service.lock.RLock()

	i := 0
	list := make([]*Model, len(service.publishers))
	for _, publisher := range service.publishers {
		list[i] = publisher.Model
		i++
	}
	return list
}
