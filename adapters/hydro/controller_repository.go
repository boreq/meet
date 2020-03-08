package hydro

import (
	"github.com/boreq/errors"
	"github.com/boreq/hydro/application/hydro"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
	"github.com/boreq/hydro/internal/eventsourcing/adapters"
	bolt "go.etcd.io/bbolt"
)

type ControllerRepository struct {
	tx         *bolt.Tx
	eventStore *eventsourcing.EventStore
}

const parentBucketName = "controllers"

func NewControllerRepository(tx *bolt.Tx) (*ControllerRepository, error) {
	persistenceAdapter := adapters.NewBoltPersistenceAdapter(tx, func(uuid eventsourcing.AggregateUUID) []adapters.BucketName {
		return []adapters.BucketName{
			[]byte(parentBucketName),
			[]byte(uuid),
			[]byte("events"),
		}
	})
	eventStore := eventsourcing.NewEventStore(controllerEventMapping, persistenceAdapter)

	return &ControllerRepository{
		tx:         tx,
		eventStore: eventStore,
	}, nil
}

func (c ControllerRepository) List() ([]*domain.Controller, error) {
	bucket := c.tx.Bucket([]byte(parentBucketName))
	if bucket == nil {
		return nil, nil
	}

	var controllers []*domain.Controller

	if err := bucket.ForEach(func(key, _ []byte) error {
		uuid, err := domain.NewControllerUUID(string(key))
		if err != nil {
			return errors.Wrap(err, "could not create a uuid")
		}

		controller, err := c.get(uuid)
		if err != nil {
			return errors.Wrapf(err, "could not get '%s'", uuid)

		}

		controllers = append(controllers, controller)
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "iteration failed")
	}

	return controllers, nil
}

func (c ControllerRepository) GetByAddress(address domain.Address) (*domain.Controller, error) {
	controllers, err := c.List()
	if err != nil {
		return nil, errors.Wrap(err, "could not list controllers")
	}

	for _, controller := range controllers {
		if controller.Address() == address {
			return controller, nil
		}
	}

	return nil, hydro.ErrControllerNotFound
}

func (c ControllerRepository) Save(controller *domain.Controller) error {
	return c.eventStore.SaveEvents(c.convertUUID(controller.UUID()), controller.PopChanges())
}

func (c ControllerRepository) get(uuid domain.ControllerUUID) (*domain.Controller, error) {
	events, err := c.eventStore.GetEvents(c.convertUUID(uuid))
	if err != nil {
		if errors.Is(err, eventsourcing.EventsNotFound) {
			return nil, hydro.ErrControllerNotFound
		}
		return nil, errors.Wrap(err, "could not get events")
	}

	return domain.NewControllerFromHistory(events)
}

func (c ControllerRepository) convertUUID(uuid domain.ControllerUUID) eventsourcing.AggregateUUID {
	return eventsourcing.AggregateUUID(uuid.String())
}
