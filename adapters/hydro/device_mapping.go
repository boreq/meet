package hydro

import (
	"encoding/json"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
)

var deviceEventMapping = eventsourcing.Mapping{
	"created_v1": eventsourcing.EventMapping{
		Marshal: func(event eventsourcing.Event) ([]byte, error) {
			e := event.(domain.DeviceCreated)

			transportEvent := deviceCreated{
				UUID:           e.UUID.String(),
				ControllerUUID: e.ControllerUUID.String(),
				ID:             e.ID.String(),
			}

			return json.Marshal(transportEvent)
		},
		Unmarshal: func(bytes []byte) (eventsourcing.Event, error) {
			var transportEvent deviceCreated

			if err := json.Unmarshal(bytes, &transportEvent); err != nil {
				return nil, errors.Wrap(err, "could not unmarshal json")
			}

			uuid, err := domain.NewDeviceUUID(transportEvent.UUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create a uuid")
			}

			controllerUUID, err := domain.NewControllerUUID(transportEvent.ControllerUUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create an address")
			}

			id, err := domain.NewDeviceID(transportEvent.ID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create an id")
			}

			return domain.DeviceCreated{
				UUID:           uuid,
				ControllerUUID: controllerUUID,
				ID:             id,
			}, nil
		},
	},
}

type deviceCreated struct {
	UUID           string `json:"uuid"`
	ControllerUUID string `json:"controllerUUID"`
	ID             string `json:"deviceID"`
}
