package hydro_test

import (
	"testing"

	"github.com/boreq/hydro/adapters/hydro"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/fixture"
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
)

func TestDeviceRepository(t *testing.T) {
	db, cleanup := fixture.Bolt(t)
	defer cleanup()

	controllerUUID := domain.MustNewControllerUUID("controller-uuid")
	address := domain.MustNewAddress("controller-address")
	controller, err := domain.NewController(controllerUUID, address)
	require.NoError(t, err)

	err = db.Update(func(tx *bbolt.Tx) error {
		r, err := hydro.NewControllerRepository(tx)
		require.NoError(t, err)

		err = r.Save(controller)
		require.NoError(t, err)

		return nil
	})
	require.NoError(t, err)

	deviceUUID := domain.MustNewDeviceUUID("device-uuid")
	deviceID := domain.MustNewDeviceID("device-id")
	device, err := domain.NewDevice(deviceUUID, controllerUUID, deviceID)
	require.NoError(t, err)

	err = db.Update(func(tx *bbolt.Tx) error {
		r, err := hydro.NewDeviceRepository(tx)
		require.NoError(t, err)

		devices, err := r.ListByController(controllerUUID)
		require.NoError(t, err)
		require.Empty(t, devices)

		err = r.Save(device)
		require.NoError(t, err)

		devices, err = r.ListByController(controllerUUID)
		require.NoError(t, err)
		require.Equal(t, []*domain.Device{device}, devices)

		return nil
	})
	require.NoError(t, err)
}
