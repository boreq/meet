package domain

import "github.com/boreq/errors"

type Controller struct {
	uuid    ControllerUUID
	devices []Device
}

func NewController(uuid ControllerUUID) (*Controller, error) {
	return &Controller{
		uuid: uuid,
	}, nil
}

func (r *Controller) AddDevice(device Device) error {
	return errors.New("not implemented")
}

func (r *Controller) RemoveDevice(deviceUUID DeviceUUID) error {
	return errors.New("not implemented")
}

func (r *Controller) Devices() []Device {
	return r.devices
}
