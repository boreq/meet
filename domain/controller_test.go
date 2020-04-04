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

	events := controller.PopChanges()

	require.Equal(t, []eventsourcing.Event{
		domain.ControllerCreated{
			UUID:    controllerUUID,
			Address: address,
		},
	}, events.Payloads())

	controllerFromHistory, err := domain.NewControllerFromHistory(events)
	require.NoError(t, err)

	require.Empty(t, controllerFromHistory.PopChanges())
	require.Equal(t, controller, controllerFromHistory)
}
