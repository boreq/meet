package hydro

import "github.com/boreq/hydro/domain"

type Controller struct {
	UUID    string `json:"uuid"`
	Address string `json:"address"`
}

type Device struct {
	UUID string `json:"uuid"`
	ID   string `json:"id"`
}

func toControllers(controllers []*domain.Controller) []Controller {
	rv := make([]Controller, 0)
	for _, controller := range controllers {
		rv = append(rv, toController(controller))
	}
	return rv
}

func toController(controller *domain.Controller) Controller {
	return Controller{
		UUID:    controller.UUID().String(),
		Address: controller.Address().String(),
	}
}

func toDevices(devices []*domain.Device) []Device {
	rv := make([]Device, 0)
	for _, device := range devices {
		rv = append(rv, toDevice(device))
	}
	return rv
}

func toDevice(device *domain.Device) Device {
	return Device{
		UUID: device.UUID().String(),
		ID:   device.ID().String(),
	}
}
