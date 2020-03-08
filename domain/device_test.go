package domain_test

import (
	"testing"

	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
	"github.com/stretchr/testify/require"
)

func TestDevice(t *testing.T) {
	deviceUUID, err := domain.NewDeviceUUID("device-uuid")
	require.NoError(t, err)

	device, err := domain.NewDevice(deviceUUID)
	require.NoError(t, err)

	require.Equal(t, deviceUUID, device.UUID())
	require.Zero(t, device.Schedule())
	require.Equal(t, domain.DeviceModeAuto, device.Mode())

	schedule := domain.MustNewSchedule(
		[]domain.Period{
			domain.MustNewPeriod(
				domain.MustNewTime(11, 00),
				domain.MustNewTime(12, 00),
			),
			domain.MustNewPeriod(
				domain.MustNewTime(13, 00),
				domain.MustNewTime(14, 00),
			),
		},
	)

	err = device.SetSchedule(schedule)
	require.NoError(t, err)

	err = device.SetSchedule(schedule)
	require.NoError(t, err)

	err = device.SetMode(domain.DeviceModeOn)
	require.NoError(t, err)

	err = device.SetMode(domain.DeviceModeOn)
	require.NoError(t, err)

	events := device.PopChanges()

	require.Equal(t, []eventsourcing.Event{
		domain.DeviceCreated{
			UUID: deviceUUID,
		},
		domain.ScheduleSet{
			Schedule: schedule,
		},
		domain.ModeSet{
			Mode: domain.DeviceModeOn,
		},
	}, events.Payloads())

	deviceFromHistory, err := domain.NewDeviceFromHistory(events)
	require.NoError(t, err)

	require.Empty(t, deviceFromHistory.PopChanges())
	require.Equal(t, device, deviceFromHistory)
}
