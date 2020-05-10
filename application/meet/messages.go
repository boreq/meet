package meet

import (
	"github.com/boreq/meet/domain"
)

type IncomingMessage interface{}

type SetNameMessage struct {
	Name domain.ParticipantName
}

type UpdateVisualisationStateMessage struct {
	State domain.VisualisationState
}

type LocalSessionDescriptionMessage struct {
	TargetParticipantUUID domain.ParticipantUUID
	SessionDescription    domain.SessionDescription
}

type LocalIceCandidateMessage struct {
	TargetParticipantUUID domain.ParticipantUUID
	IceCandidate          domain.IceCandidate
}
