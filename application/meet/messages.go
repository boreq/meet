package meet

import (
	"github.com/boreq/meet/domain"
)

type IncomingMessage interface{}

type SetNameMessage struct {
	Name domain.ParticipantName
}

type BrowserSessionDescription struct {
	SessionDescription domain.RemoteDescription
}
