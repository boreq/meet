package hydro

import (
	"encoding/json"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
)

var controllerEventMapping = eventsourcing.Mapping{
	"created_v1": eventsourcing.EventMapping{
		Marshal: func(event eventsourcing.Event) ([]byte, error) {
			e := event.(domain.ControllerCreated)

			transportEvent := controllerCreated{
				UUID:    e.UUID.String(),
				Address: e.Address.String(),
			}

			return json.Marshal(transportEvent)
		},
		Unmarshal: func(bytes []byte) (eventsourcing.Event, error) {
			var transportEvent controllerCreated

			if err := json.Unmarshal(bytes, &transportEvent); err != nil {
				return nil, errors.Wrap(err, "could not unmarshal json")
			}

			uuid, err := domain.NewControllerUUID(transportEvent.UUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create a uuid")
			}

			address, err := domain.NewAddress(transportEvent.Address)
			if err != nil {
				return nil, errors.Wrap(err, "could not create an address")
			}

			return domain.ControllerCreated{
				UUID:    uuid,
				Address: address,
			}, nil
		},
	},
}

type controllerCreated struct {
	UUID    string `json:"uuid"`
	Address string `json:"address"`
}
