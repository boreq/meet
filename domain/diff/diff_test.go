package diff_test

import (
	"testing"

	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/domain/diff"
	"github.com/stretchr/testify/require"
)

func TestDevices(t *testing.T) {
	controllerUUID := domain.MustNewControllerUUID("controller-uuid")

	device1, err := domain.NewDevice(
		domain.MustNewDeviceUUID("device-uuid-1"),
		controllerUUID,
		domain.MustNewDeviceID("device-id-1"),
	)
	require.NoError(t, err)

	device2, err := domain.NewDevice(
		domain.MustNewDeviceUUID("device-uuid-2"),
		controllerUUID,
		domain.MustNewDeviceID("device-id-2"),
	)
	require.NoError(t, err)

	existingDevices := []*domain.Device{
		device1,
		device2,
	}

	notExistingDeviceId := domain.MustNewDeviceID("device-id-3")

	devices := []domain.DeviceID{
		device2.ID(),
		notExistingDeviceId,
		notExistingDeviceId,
	}

	toAdd, toRemove := diff.Devices(existingDevices, devices)

	require.Equal(t, []domain.DeviceID{
		notExistingDeviceId,
	}, toAdd)

	require.Equal(t, []*domain.Device{
		device1,
	}, toRemove)
}
