package hydro

import (
	"fmt"
	"github.com/boreq/errors"
	"github.com/boreq/hydro/application/hydro"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
	"github.com/boreq/hydro/internal/eventsourcing/adapters"
	bolt "go.etcd.io/bbolt"
)

const devicesBucket = "devices"
const devicesMappingBucket = "devices_mapping"

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

	bucket = bucket.Bucket([]byte(devicesBucket))
	if bucket == nil {
		return nil, nil
	}

	var devices []*domain.Device

	if err := bucket.ForEach(func(key, value []byte) error {
		deviceUUID, err := domain.NewDeviceUUID(string(key))
		if err != nil {
			return errors.Wrap(err, "could not create a device uuid from the key")
		}

		fmt.Println(string(key))

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

func (r *DeviceRepository) Save(device *domain.Device) error {
	if err := r.putControllerMapping(device.ControllerUUID(), device.UUID()); err != nil {
		return errors.Wrap(err, "could not put the controller mapping")
	}

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
	return eventsourcing.NewEventStore(deviceEventMapping, persistenceAdapter)
}

func (r *DeviceRepository) convertUUID(uuid domain.DeviceUUID) eventsourcing.AggregateUUID {
	return eventsourcing.AggregateUUID(uuid.String())
}

func (r *DeviceRepository) putControllerMapping(controllerUUID domain.ControllerUUID, deviceUUID domain.DeviceUUID) error {
	bucket, err := createBucket(r.tx, [][]byte{
		[]byte(devicesMappingBucket),
	})
	if err != nil {
		return errors.Wrap(err, "could not create a bucket")
	}

	if err := bucket.Put([]byte(deviceUUID.String()), []byte(controllerUUID.String())); err != nil {
		return errors.Wrap(err, "could not put the mapping entry")
	}

	return nil
}

func (r *DeviceRepository) removeControllerMapping(deviceUUID domain.DeviceUUID) error {
	bucket, err := createBucket(r.tx, [][]byte{
		[]byte(devicesMappingBucket),
	})
	if err != nil {
		return errors.Wrap(err, "could not create a bucket")
	}

	if err := bucket.Delete([]byte(deviceUUID.String())); err != nil {
		return errors.Wrap(err, "could not put the mapping entry")
	}

	return nil
}

func (r *DeviceRepository) getControllerMapping(deviceUUID domain.DeviceUUID) (domain.ControllerUUID, error) {
	bucket := getBucket(r.tx, [][]byte{
		[]byte(devicesMappingBucket),
	})
	if bucket == nil {
		return domain.ControllerUUID{}, errors.Wrap(hydro.ErrDeviceNotFound, "mapping does not exist")
	}

	value := bucket.Get([]byte(deviceUUID.String()))
	if value == nil {
		return domain.ControllerUUID{}, errors.New("invalid mapping, nil value")
	}

	return domain.NewControllerUUID(string(value))
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

func createBucket(tx *bolt.Tx, bucketNames [][]byte) (bucket *bolt.Bucket, err error) {
	bucket, err = tx.CreateBucketIfNotExists(bucketNames[0])
	if err != nil {
		return nil, errors.Wrap(err, "could not create a bucket")
	}

	for i := 1; i < len(bucketNames); i++ {
		bucket, err = bucket.CreateBucketIfNotExists(bucketNames[i])
		if err != nil {
			return nil, errors.Wrap(err, "could not create a bucket")
		}
	}

	return bucket, nil
}
