package hydro_test

import (
	"context"
	"testing"

	"github.com/boreq/hydro/application/hydro"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
	"github.com/boreq/hydro/internal/wire"
	"github.com/stretchr/testify/require"
)

func TestSetControllerDevices(t *testing.T) {
	ctx := context.Background()

	app, err := wire.BuildUnitTestHydroApplication()
	require.NoError(t, err)

	controllerUUID := domain.MustNewControllerUUID("controller-uuid")
	controllerAddress := domain.MustNewAddress("controller-address")

	controller, err := domain.NewController(controllerUUID, controllerAddress)
	require.NoError(t, err)

	device1, err := domain.NewDevice(
		domain.MustNewDeviceUUID("device1-uuid"),
		controllerUUID,
		domain.MustNewDeviceID("device1-id"),
	)
	require.NoError(t, err)

	device2, err := domain.NewDevice(
		domain.MustNewDeviceUUID("device2-uuid"),
		controllerUUID,
		domain.MustNewDeviceID("device2-id"),
	)
	require.NoError(t, err)

	device3, err := domain.NewDevice(
		domain.MustNewDeviceUUID("device3-uuid"),
		controllerUUID,
		domain.MustNewDeviceID("device3-id"),
	)
	require.NoError(t, err)

	err = controller.AddDevice(device1.UUID())
	require.NoError(t, err)

	err = controller.AddDevice(device2.UUID())
	require.NoError(t, err)

	err = controller.AddDevice(device3.UUID())
	require.NoError(t, err)

	err = app.Repositories.Device.Save(device1)
	require.NoError(t, err)

	err = app.Repositories.Device.Save(device2)
	require.NoError(t, err)

	err = app.Repositories.Device.Save(device3)
	require.NoError(t, err)

	err = app.Repositories.Controller.Save(controller)
	require.NoError(t, err)

	deviceID4 := domain.MustNewDeviceID("device4")
	deviceID5 := domain.MustNewDeviceID("device5")

	cmd := hydro.SetControllerDevices{
		ControllerUUID: controllerUUID,
		Devices: []domain.DeviceID{
			device3.ID(),
			deviceID4,
			deviceID5,
		},
	}

	err = app.Hydro.SetControllerDevices.Execute(ctx, cmd)
	require.NoError(t, err)

	expectedDevice4UUID := domain.MustNewDeviceUUID("uuid-1")
	expectedDevice5UUID := domain.MustNewDeviceUUID("uuid-2")

	require.Equal(t,
		[]eventsourcing.Event{
			domain.ControllerCreated{
				UUID:    controllerUUID,
				Address: controllerAddress,
			},
			domain.DeviceAdded{
				DeviceUUID: device1.UUID(),
			},
			domain.DeviceAdded{
				DeviceUUID: device2.UUID(),
			},
			domain.DeviceAdded{
				DeviceUUID: device3.UUID(),
			},
			domain.DeviceAdded{
				DeviceUUID: expectedDevice4UUID,
			},
			domain.DeviceAdded{
				DeviceUUID: expectedDevice5UUID,
			},
			domain.DeviceRemoved{
				DeviceUUID: device1.UUID(),
			},
			domain.DeviceRemoved{
				DeviceUUID: device2.UUID(),
			},
		},
		app.Repositories.Controller.Events[controllerUUID].Payloads())

	require.Empty(t, app.Repositories.Device.Events[controllerUUID][device1.UUID()].Payloads())
	require.Empty(t, app.Repositories.Device.Events[controllerUUID][device2.UUID()].Payloads())

	require.Equal(t,
		[]eventsourcing.Event{
			domain.DeviceCreated{
				UUID:           device3.UUID(),
				ControllerUUID: controllerUUID,
				ID:             device3.ID(),
			},
		},
		app.Repositories.Device.Events[controllerUUID][device3.UUID()].Payloads())

	require.Equal(t,
		[]eventsourcing.Event{
			domain.DeviceCreated{
				UUID:           expectedDevice4UUID,
				ControllerUUID: controllerUUID,
				ID:             deviceID4,
			},
		},
		app.Repositories.Device.Events[controllerUUID][expectedDevice4UUID].Payloads())

	require.Equal(t,
		[]eventsourcing.Event{
			domain.DeviceCreated{
				UUID:           expectedDevice5UUID,
				ControllerUUID: controllerUUID,
				ID:             deviceID5,
			},
		},
		app.Repositories.Device.Events[controllerUUID][expectedDevice5UUID].Payloads())
}
