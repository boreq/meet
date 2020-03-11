package hydro

import (
	"github.com/boreq/errors"
	"github.com/boreq/hydro/application/hydro"
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

func (r *DeviceRepository) ListByController(controllerUUID domain.ControllerUUID) ([]*domain.Device, error) {
	bucket := getBucket(r.tx,
		[][]byte{
			[]byte(controllersBucket),
			[]byte(controllerUUID.String()),
		},
	)

	if bucket == nil {
		return nil, hydro.ErrControllerNotFound
	}

	var devices []*domain.Device

	if err := bucket.ForEach(func(key, value []byte) error {
		deviceUUID, err := domain.NewDeviceUUID(string(key))
		if err != nil {
			return errors.Wrap(err, "could not create a device uuid from the key")
		}
		device, err := r.get(controllerUUID, deviceUUID)
		if err != nil {
			return errors.Wrap(err, "could not get a device")
		}

		devices = append(devices, device)

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "foreach failed")
	}

	return devices, nil
}

func (r *DeviceRepository) Remove(uuid domain.DeviceUUID) error {
	return errors.New("not implemented")
}

func (r *DeviceRepository) Save(device *domain.Device) error {
	return r.getEventStore(device.ControllerUUID()).
		SaveEvents(r.convertUUID(device.UUID()), device.PopChanges())
}

func (r *DeviceRepository) get(controllerUUID domain.ControllerUUID, deviceUUID domain.DeviceUUID) (*domain.Device, error) {
	events, err := r.getEventStore(controllerUUID).
		GetEvents(r.convertUUID(deviceUUID))
	if err != nil {
		//if errors.Is(err, eventsourcing.EventsNotFound) {
		//	return nil, hydro.ErrControllerNotFound
		//}
		return nil, errors.Wrap(err, "could not get events")
	}

	return domain.NewDeviceFromHistory(events)
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

func getBucket(tx *bolt.Tx, bucketNames [][]byte) *bolt.Bucket {
	bucket := tx.Bucket(bucketNames[0])

	if bucket == nil {
		return nil
	}

	for i := 1; i < len(bucketNames); i++ {
		bucket = bucket.Bucket(bucketNames[i])
		if bucket == nil {
			return nil
		}
	}

	return bucket
}
