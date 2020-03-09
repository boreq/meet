package domain

import (
	"fmt"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/internal/eventsourcing"
)

type Device struct {
	uuid           DeviceUUID
	controllerUUID ControllerUUID
	id             DeviceID
	schedule       Schedule
	mode           DeviceMode

	es eventsourcing.EventSourcing
}

func NewDevice(uuid DeviceUUID, controllerUUID ControllerUUID, id DeviceID) (*Device, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	if controllerUUID.IsZero() {
		return nil, errors.New("zero value of controller uuid")
	}

	if id.IsZero() {
		return nil, errors.New("zero value of id")
	}

	device := &Device{}

	event := DeviceCreated{uuid, controllerUUID, id}
	if err := device.update(event); err != nil {
		return nil, errors.Wrap(err, "could not consume an event")
	}

	return device, nil
}

func NewDeviceFromHistory(events []eventsourcing.EventSourcingEvent) (*Device, error) {
	device := &Device{}

	for _, event := range events {
		if err := device.update(event.Event); err != nil {
			return nil, errors.Wrap(err, "could not process an event")
		}
		device.es.LoadVersion(event)
	}

	device.es.PopChanges()

	return device, nil
}

func (d *Device) SetSchedule(schedule Schedule) error {
	// idempotence
	if d.schedule.Equal(schedule) {
		return nil
	}

	event := ScheduleSet{schedule}
	return d.update(event)
}

func (d *Device) SetMode(mode DeviceMode) error {
	if mode.IsZero() {
		return errors.New("zero value of device mode")
	}

	// idempotence
	if d.mode == mode {
		return nil
	}

	event := ModeSet{mode}
	return d.update(event)
}

func (d *Device) UUID() DeviceUUID {
	return d.uuid
}

func (d *Device) ControllerUUID() ControllerUUID {
	return d.controllerUUID
}

func (d *Device) ID() DeviceID {
	return d.id
}

func (d *Device) Schedule() Schedule {
	return d.schedule
}

func (d *Device) Mode() DeviceMode {
	return d.mode
}

func (d *Device) PopChanges() eventsourcing.EventSourcingEvents {
	return d.es.PopChanges()
}

func (d *Device) update(event eventsourcing.Event) error {
	switch e := event.(type) {
	case DeviceCreated:
		d.handleDeviceCreated(e)
	case ScheduleSet:
		d.handleScheduleSet(e)
	case ModeSet:
		d.handleModeSet(e)
	default:
		return fmt.Errorf("unsupported event '%T'", event)
	}

	return d.es.Record(event)
}

func (d *Device) handleDeviceCreated(e DeviceCreated) {
	d.uuid = e.UUID
	d.controllerUUID = e.ControllerUUID
	d.id = e.ID
	d.schedule = Schedule{}
	d.mode = DeviceModeAuto
}

func (d *Device) handleScheduleSet(e ScheduleSet) {
	d.schedule = e.Schedule
}

func (d *Device) handleModeSet(e ModeSet) {
	d.mode = e.Mode
}

type DeviceCreated struct {
	UUID           DeviceUUID
	ControllerUUID ControllerUUID
	ID             DeviceID
}

func (c DeviceCreated) EventType() eventsourcing.EventType {
	return "created_v1"
}

type ScheduleSet struct {
	Schedule Schedule
}

func (c ScheduleSet) EventType() eventsourcing.EventType {
	return "schedule_set_v1"
}

type ModeSet struct {
	Mode DeviceMode
}

func (c ModeSet) EventType() eventsourcing.EventType {
	return "mode_set_v1"
}

type DeviceMode struct {
	v string
}

func (m DeviceMode) IsZero() bool {
	return m.v == ""
}

var (
	DeviceModeAuto = DeviceMode{"auto"}
	DeviceModeOn   = DeviceMode{"on"}
	DeviceModeOff  = DeviceMode{"off"}
)

type DeviceID struct {
	id string
}

func NewDeviceID(id string) (DeviceID, error) {
	if id == "" {
		return DeviceID{}, errors.New("address can not be empty")
	}

	return DeviceID{
		id: id,
	}, nil
}

func MustNewDeviceID(id string) DeviceID {
	v, err := NewDeviceID(id)
	if err != nil {
		panic(err)
	}
	return v
}

func (i DeviceID) String() string {
	return i.id
}

func (i DeviceID) IsZero() bool {
	return i == DeviceID{}
}
