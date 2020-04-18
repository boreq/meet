package meet

import (
	"github.com/boreq/meet/domain"
)

type IncomingMessage interface{}

type SetNameMessage struct {
	Name domain.ParticipantName
}

type PassLocalSessionDescription struct {
	TargetParticipant  domain.ParticipantUUID
	SessionDescription domain.SessionDescription
}
