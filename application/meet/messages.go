package meet

import (
	"github.com/boreq/meet/domain"
)

type IncomingMessage interface{}

type SetNameMessage struct {
	Name domain.ParticipantName
}

type LocalSessionDescription struct {
	TargetParticipantUUID domain.ParticipantUUID
	SessionDescription    domain.SessionDescription
}

type LocalIceCandidate struct {
	TargetParticipantUUID domain.ParticipantUUID
	IceCandidate          domain.IceCandidate
}
