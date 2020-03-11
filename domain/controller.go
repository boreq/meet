package domain

import (
	"fmt"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/internal/eventsourcing"
)

type Controller struct {
	uuid    ControllerUUID
	address Address
	devices []DeviceUUID

	es eventsourcing.EventSourcing
}

func NewController(uuid ControllerUUID, address Address) (*Controller, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	if address.IsZero() {
		return nil, errors.New("zero value of address")
	}

	controller := &Controller{}

	event := ControllerCreated{uuid, address}
	if err := controller.update(event); err != nil {
		return nil, errors.Wrap(err, "could not consume an event")
	}

	return controller, nil
}

func NewControllerFromHistory(events []eventsourcing.EventSourcingEvent) (*Controller, error) {
	controller := &Controller{}

	for _, event := range events {
		if err := controller.update(event.Event); err != nil {
			return nil, errors.Wrap(err, "could not process an event")
		}
		controller.es.LoadVersion(event)
	}

	controller.es.PopChanges()

	return controller, nil
}

func (c *Controller) AddDevice(deviceUUID DeviceUUID) error {
	if deviceUUID.IsZero() {
		return errors.New("zero value of device uuid")
	}

	if c.hasDevice(deviceUUID) {
		return errors.New("this device already exists")
	}

	return c.update(DeviceAdded{deviceUUID})
}

func (c *Controller) RemoveDevice(deviceUUID DeviceUUID) error {
	if deviceUUID.IsZero() {
		return errors.New("zero value of device uuid")
	}

	if !c.hasDevice(deviceUUID) {
		return errors.New("this device does not exist")
	}

	return c.update(DeviceRemoved{deviceUUID})
}

func (c *Controller) hasDevice(deviceUUID DeviceUUID) bool {
	for _, device := range c.devices {
		if device == deviceUUID {
			return true
		}
	}
	return false
}

func (c *Controller) UUID() ControllerUUID {
	return c.uuid
}

func (c *Controller) Address() Address {
	return c.address
}

func (c *Controller) Devices() []DeviceUUID {
	tmp := make([]DeviceUUID, len(c.devices))
	copy(tmp, c.devices)
	return tmp
}

func (c *Controller) PopChanges() eventsourcing.EventSourcingEvents {
	return c.es.PopChanges()
}

func (c *Controller) update(event eventsourcing.Event) error {
	switch e := event.(type) {
	case ControllerCreated:
		c.handleCreated(e)
	case DeviceAdded:
		c.handleDeviceAdded(e)
	case DeviceRemoved:
		c.handleDeviceRemoved(e)
	default:
		return fmt.Errorf("unsupported event '%T'", event)
	}

	return c.es.Record(event)
}

func (c *Controller) handleCreated(e ControllerCreated) {
	c.uuid = e.UUID
	c.address = e.Address
}

func (c *Controller) handleDeviceAdded(e DeviceAdded) {
	c.devices = append(c.devices, e.DeviceUUID)
}

func (c *Controller) handleDeviceRemoved(e DeviceRemoved) {
	for i := len(c.devices) - 1; i >= 0; i-- {
		if c.devices[i] == e.DeviceUUID {
			c.devices = append(c.devices[:i], c.devices[i+1:]...)
		}
	}
}

type ControllerCreated struct {
	UUID    ControllerUUID
	Address Address
}

func (c ControllerCreated) EventType() eventsourcing.EventType {
	return "created_v1"
}

type DeviceAdded struct {
	DeviceUUID DeviceUUID
}

func (c DeviceAdded) EventType() eventsourcing.EventType {
	return "device_added_v1"
}

type DeviceRemoved struct {
	DeviceUUID DeviceUUID
}

func (c DeviceRemoved) EventType() eventsourcing.EventType {
	return "device_removed_v1"
}
