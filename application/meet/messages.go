package meet

import (
	"github.com/boreq/meet/domain"
	"github.com/boreq/meet/domain/webrtc"
)

type IncomingMessage interface{}

type SetNameMessage struct {
	Name domain.Name
}

type BrowserSessionDescription struct {
	SessionDescription webrtc.RemoteDescription
}
