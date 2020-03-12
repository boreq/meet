package hydro

import (
	"context"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/domain/diff"
)

type SetControllerDevices struct {
	ControllerUUID domain.ControllerUUID
	Devices        []domain.DeviceID
}

type SetControllerDevicesHandler struct {
	transactionProvider TransactionProvider
	uuidGenerator       UUIDGenerator
}

func NewSetControllerDevicesHandler(
	transactionProvider TransactionProvider,
	uuidGenerator UUIDGenerator,
) *SetControllerDevicesHandler {
	return &SetControllerDevicesHandler{
		transactionProvider: transactionProvider,
		uuidGenerator:       uuidGenerator,
	}
}

func (h *SetControllerDevicesHandler) Execute(ctx context.Context, cmd SetControllerDevices) error {
	return h.transactionProvider.Transact(ctx, func(a *TransactableAdapters) error {
		controller, err := a.Controllers.Get(cmd.ControllerUUID)
		if err != nil {
			return errors.Wrap(err, "could not get the controller")
		}

		devices, err := a.Devices.ListByController(controller.UUID())
		if err != nil {
			return errors.Wrap(err, "could not list the devices")
		}

		toAdd, toRemove := diff.Devices(devices, cmd.Devices)

		for _, deviceId := range toAdd {
			device, err := h.newDevice(controller.UUID(), deviceId)
			if err != nil {
				return errors.Wrap(err, "could not create a device")
			}

			if err := controller.AddDevice(device.UUID()); err != nil {
				return errors.Wrap(err, "could not add a device to the controller")
			}

			if err := a.Devices.Save(device); err != nil {
				return errors.Wrap(err, "could not save the device")
			}
		}

		for _, device := range toRemove {
			if err := controller.RemoveDevice(device.UUID()); err != nil {
				return errors.Wrap(err, "could not remove a device from the controller")
			}

			// todo
			//if err := a.Devices.Remove(device.UUID()); err != nil {
			//	return errors.Wrap(err, "could not remove a device")
			//}
		}

		return a.Controllers.Save(controller)
	})
}

func (h *SetControllerDevicesHandler) newDevice(controllerUUID domain.ControllerUUID, deviceID domain.DeviceID) (*domain.Device, error) {
	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate a uuid")
	}

	deviceUUID, err := domain.NewDeviceUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a device uuid")
	}

	return domain.NewDevice(deviceUUID, controllerUUID, deviceID)
}
