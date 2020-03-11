package hydro

import (
	"sort"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
)

type DeviceRepositoryMock struct {
	Events map[domain.ControllerUUID]map[domain.DeviceUUID]eventsourcing.EventSourcingEvents
}

func NewDeviceRepositoryMock() *DeviceRepositoryMock {
	return &DeviceRepositoryMock{
		Events: make(map[domain.ControllerUUID]map[domain.DeviceUUID]eventsourcing.EventSourcingEvents),
	}
}

func (d DeviceRepositoryMock) ListByController(uuid domain.ControllerUUID) ([]*domain.Device, error) {
	var devices []*domain.Device
	for _, events := range d.Events[uuid] {
		device, err := domain.NewDeviceFromHistory(events)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a device from history")
		}

		devices = append(devices, device)
	}

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].UUID().String() < devices[j].UUID().String()
	})

	return devices, nil
}

func (d DeviceRepositoryMock) Remove(uuid domain.DeviceUUID) error {
	for _, subMap := range d.Events {
		delete(subMap, uuid)
	}
	return nil
}

func (d DeviceRepositoryMock) Save(device *domain.Device) error {
	if _, ok := d.Events[device.ControllerUUID()]; !ok {
		d.Events[device.ControllerUUID()] = make(map[domain.DeviceUUID]eventsourcing.EventSourcingEvents)
	}
	d.Events[device.ControllerUUID()][device.UUID()] = append(d.Events[device.ControllerUUID()][device.UUID()], device.PopChanges()...)
	return nil
}
