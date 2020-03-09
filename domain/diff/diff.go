package diff

import (
	"github.com/boreq/hydro/domain"
)

func Devices(existingDevices []*domain.Device, devices []domain.DeviceID) (toAdd []domain.DeviceID, toRemove []*domain.Device) {
	for _, existingDevice := range existingDevices {
		if !hasDevice(devices, existingDevice.ID()) {
			if !hasExistingDevice(toRemove, existingDevice.ID()) {
				toRemove = append(toRemove, existingDevice)
			}
		}
	}

	for _, device := range devices {
		if !hasExistingDevice(existingDevices, device) {
			if !hasDevice(toAdd, device) {
				toAdd = append(toAdd, device)
			}
		}
	}

	return
}

func hasExistingDevice(existingDevices []*domain.Device, deviceID domain.DeviceID) bool {
	for _, existingDevice := range existingDevices {
		if existingDevice.ID() == deviceID {
			return true
		}
	}
	return false
}

func hasDevice(devices []domain.DeviceID, deviceID domain.DeviceID) bool {
	for _, device := range devices {
		if device == deviceID {
			return true
		}
	}
	return false
}
