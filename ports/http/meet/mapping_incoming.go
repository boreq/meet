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
	UpdateVisualisationStateMessage: func(payload []byte) (meet.IncomingMessage, error) {
		var transport UpdateVisualisationStateMsg
		if err := json.Unmarshal(payload, &transport); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}

		state, err := domain.NewVisualisationState(transport.State)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a visualisation state")
		}

		return meet.UpdateVisualisationStateMessage{
			State: state,
		}, nil
	},
	LocalSessionDescriptionMessage: func(payload []byte) (meet.IncomingMessage, error) {
		var transport LocalSessionDescriptionMsg
		if err := json.Unmarshal(payload, &transport); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}

		targetParticipantUUID, err := domain.NewParticipantUUID(transport.TargetParticipantUUID)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a target participant uuid")
		}

		sessionDescription, err := domain.NewSessionDescription(transport.SessionDescription)
		if err != nil {
			return nil, errors.Wrap(err, "could not create an ice candidate")
		}

		return meet.LocalSessionDescriptionMessage{
			TargetParticipantUUID: targetParticipantUUID,
			SessionDescription:    sessionDescription,
		}, nil
	},
	LocalIceCandidateMessage: func(payload []byte) (meet.IncomingMessage, error) {
		var transport LocalIceCandidateMsg
		if err := json.Unmarshal(payload, &transport); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}

		targetParticipantUUID, err := domain.NewParticipantUUID(transport.TargetParticipantUUID)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a target participant uuid")
		}

		iceCandidate, err := domain.NewIceCandidate(transport.IceCandidate)
		if err != nil {
			return nil, errors.Wrap(err, "could not create an ice candidate")
		}

		return meet.LocalIceCandidateMessage{
			TargetParticipantUUID: targetParticipantUUID,
			IceCandidate:          iceCandidate,
		}, nil
	},
}

type SetNameMsg struct {
	Name string `json:"name"`
}

type UpdateVisualisationStateMsg struct {
	State string `json:"state"`
}

type LocalSessionDescriptionMsg struct {
	TargetParticipantUUID string `json:"targetParticipantUUID"`
	SessionDescription    string `json:"sessionDescription"`
}

type LocalIceCandidateMsg struct {
	TargetParticipantUUID string `json:"targetParticipantUUID"`
	IceCandidate          string `json:"iceCandidate"`
}
