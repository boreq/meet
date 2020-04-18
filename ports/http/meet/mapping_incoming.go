package meet

import (
	"encoding/json"

	"github.com/boreq/errors"
	"github.com/boreq/meet/application/meet"
	"github.com/boreq/meet/domain"
)

var incomingMapping = map[IncomingMessageType]func(payload []byte) (meet.IncomingMessage, error){
	SetNameMessage: func(payload []byte) (meet.IncomingMessage, error) {
		var transport SetNameMsg
		if err := json.Unmarshal(payload, &transport); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}

		name, err := domain.NewParticipantName(transport.Name)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a name")
		}

		return meet.SetNameMessage{
			Name: name,
		}, nil
	},
}

type SetNameMsg struct {
	Name string `json:"name"`
}
