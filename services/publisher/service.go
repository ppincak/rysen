package publisher

import (
	"strings"
	"sync"

	"github.com/ppincak/rysen/api"

	"github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/pub"
)

type Service struct {
	bus          *bus.Bus
	kafkaBrokers []string

	publishers map[string]*KafkaPublisher
	lock       *sync.RWMutex
}

func NewService(bus *bus.Bus, kafkaBrokers []string) *Service {
	return &Service{
		bus:          bus,
		kafkaBrokers: kafkaBrokers,
	}
}

func (service *Service) assembleKey(model *Model) string {
	return strings.Join([]string{model.ReadTopic, model.KafkaTopic}, "/")
}

func (service *Service) CreatePublisher(model *Model) (pub.Publisher, error) {
	defer service.lock.Unlock()
	service.lock.Unlock()

	key := service.assembleKey(model)

	if _, ok := service.publishers[key]; ok {
		return nil, api.NewError("Publisher [%s] already exists", key)
	}
	publisher := NewKafkaPublisher(model, nil)
	service.publishers[key] = publisher

	return publisher, nil
}

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
