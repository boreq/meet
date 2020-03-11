package domain_test

import (
	"testing"

	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
	"github.com/stretchr/testify/require"
)

func TestController(t *testing.T) {
	controllerUUID, err := domain.NewControllerUUID("controller-uuid")
	require.NoError(t, err)

	address, err := domain.NewAddress("controller-address")
	require.NoError(t, err)

	controller, err := domain.NewController(controllerUUID, address)
	require.NoError(t, err)

	require.Equal(t, controllerUUID, controller.UUID())
	require.Equal(t, address, controller.Address())
	require.Empty(t, controller.Devices())

	deviceUUID1, err := domain.NewDeviceUUID("device-uuid-1")
	require.NoError(t, err)

	deviceUUID2, err := domain.NewDeviceUUID("device-uuid-2")
	require.NoError(t, err)

	err = controller.AddDevice(deviceUUID1)
	require.NoError(t, err)

	err = controller.AddDevice(deviceUUID1)
	require.Error(t, err)

	err = controller.AddDevice(deviceUUID2)
	require.NoError(t, err)

	err = controller.AddDevice(deviceUUID2)
	require.Error(t, err)

	err = controller.RemoveDevice(deviceUUID1)
	require.NoError(t, err)

	err = controller.RemoveDevice(deviceUUID1)
	require.Error(t, err)

	require.Equal(t, []domain.DeviceUUID{
		deviceUUID2,
	}, controller.Devices())

	events := controller.PopChanges()

	require.Equal(t, []eventsourcing.Event{
		domain.ControllerCreated{
			UUID:    controllerUUID,
			Address: address,
		},
		domain.DeviceAdded{
			DeviceUUID: deviceUUID1,
		},
		domain.DeviceAdded{
			DeviceUUID: deviceUUID2,
		},
		domain.DeviceRemoved{
			DeviceUUID: deviceUUID1,
		},
	}, events.Payloads())

	controllerFromHistory, err := domain.NewControllerFromHistory(events)
	require.NoError(t, err)

	require.Empty(t, controllerFromHistory.PopChanges())
	require.Equal(t, controller, controllerFromHistory)
}
