package hydro

import (
	"errors"

	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
	"github.com/boreq/hydro/internal/eventsourcing/adapters"
	bolt "go.etcd.io/bbolt"
)

const devicesBucket = "devices"

type DeviceRepository struct {
	tx *bolt.Tx
}

func NewDeviceRepository(tx *bolt.Tx) (*DeviceRepository, error) {
	return &DeviceRepository{
		tx: tx,
	}, nil
}

func (r *DeviceRepository) ListByController(uuid domain.ControllerUUID) ([]*domain.Device, error) {
	return nil, errors.New("not implemented")
}

func (r *DeviceRepository) Remove(uuid domain.DeviceUUID) error {
	return errors.New("not implemented")
}

func (r *DeviceRepository) Save(device *domain.Device) error {
	return r.getEventStore(device.ControllerUUID()).
		SaveEvents(r.convertUUID(device.UUID()), device.PopChanges())
}

func (r *DeviceRepository) getEventStore(controllerUUID domain.ControllerUUID) *eventsourcing.EventStore {
	persistenceAdapter := adapters.NewBoltPersistenceAdapter(r.tx, func(uuid eventsourcing.AggregateUUID) []adapters.BucketName {
		return []adapters.BucketName{
			[]byte(controllersBucket),
			[]byte(controllerUUID.String()),
			[]byte(devicesBucket),
			[]byte(uuid),
			[]byte("events"),
		}
	})
	return eventsourcing.NewEventStore(controllerEventMapping, persistenceAdapter)
}

func (r *DeviceRepository) convertUUID(uuid domain.DeviceUUID) eventsourcing.AggregateUUID {
	return eventsourcing.AggregateUUID(uuid.String())
}
