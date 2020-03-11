package scanner

import (
	"context"
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/adapters/hydro/controller"
	"github.com/boreq/hydro/application"
	"github.com/boreq/hydro/application/hydro"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/logging"
)

type ControllerClient interface {
	GetState(ctx context.Context, address domain.Address) (controller.ControllerState, error)
}

type Scanner struct {
	client    ControllerClient
	app       *application.Application
	scanEvery time.Duration
	log       logging.Logger
}

func NewScanner(client ControllerClient, app *application.Application, scanEvery time.Duration) *Scanner {
	return &Scanner{
		client:    client,
		app:       app,
		scanEvery: scanEvery,
		log:       logging.New("ports/scanner"),
	}
}

func (s *Scanner) Run(ctx context.Context) error {
	for {
		if err := s.run(ctx); err != nil {
			s.log.Error("scanner error", "err", err)
		}

		select {
		case <-time.After(s.scanEvery):
			continue
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Scanner) run(ctx context.Context) error {
	controllers, err := s.app.Hydro.ListControllersHandler.Execute(ctx)
	if err != nil {
		return errors.Wrap(err, "error listing controllers")
	}

	controllerStates, err := s.query(ctx, controllers)
	if err != nil {
		return errors.Wrap(err, "error querying controllers")
	}

	for controllerUUID, state := range controllerStates {
		deviceIds, err := toDeviceIds(state)
		if err != nil {
			return errors.Wrap(err, "could not convert the state to device ids")
		}

		cmd := hydro.SetControllerDevices{
			ControllerUUID: controllerUUID,
			Devices:        deviceIds,
		}

		if err := s.app.Hydro.SetControllerDevicesHandler.Execute(ctx, cmd); err != nil {
			return errors.Wrapf(err, "could not set controller devices uuid='%s'", controllerUUID)
		}
	}

	return nil
}

func (s *Scanner) query(ctx context.Context, controllers []*domain.Controller) (map[domain.ControllerUUID]controller.ControllerState, error) {
	c := make(chan queryResult)
	for i := range controllers {
		go func(i int) {
			s.log.Debug("querying a controller", "address", controllers[i].Address())
			state, err := s.client.GetState(ctx, controllers[i].Address())
			c <- newQueryResult(controllers[i].UUID(), state, err)
		}(i)
	}

	states := make(map[domain.ControllerUUID]controller.ControllerState)

	for i := 0; i < len(controllers); i++ {
		select {
		case result := <-c:
			uuid, state, err := result.Unpack()
			if err != nil {
				s.log.Error("error querying a controller", "err", err)
				continue
			}
			states[uuid] = state
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	s.log.Debug("all controllers have replied")

	return states, nil
}

func toDeviceIds(state controller.ControllerState) ([]domain.DeviceID, error) {
	var deviceIds []domain.DeviceID

	for _, deviceState := range state.Devices {
		deviceId, err := domain.NewDeviceID(deviceState.Id)
		if err != nil {
			return nil, errors.Wrap(err, "error creating a device id")
		}
		deviceIds = append(deviceIds, deviceId)
	}

	return deviceIds, nil
}

type queryResult struct {
	uuid  domain.ControllerUUID
	state controller.ControllerState
	err   error
}

func newQueryResult(uuid domain.ControllerUUID, state controller.ControllerState, err error) queryResult {
	return queryResult{uuid: uuid, state: state, err: err}
}

func (r queryResult) Unpack() (domain.ControllerUUID, controller.ControllerState, error) {
	return r.uuid, r.state, r.err
}
