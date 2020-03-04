package domain

import (
	"github.com/boreq/errors"
)

type Device struct {
	uuid       DeviceUUID
	schedule   Schedule
}

func NewDevice(uuid DeviceUUID, schedule Schedule) (Device, error) {
	if uuid.IsZero() {
		return Device{}, errors.New("zero value of uuid")
	}

	return Device{
		uuid:       uuid,
		schedule:   schedule,
	}, nil
}

func (d Device) SetSchedule(schedule Schedule) error {
	d.schedule = schedule
	return nil
}

func (d Device) UUID() DeviceUUID {
	return d.uuid
}

func (d Device) Schedule() Schedule {
	return d.schedule
}
