package webrtc

import (
	"github.com/boreq/errors"
	"github.com/boreq/meet/adapters/meet"
	"github.com/boreq/meet/domain"
	"github.com/boreq/meet/internal/logging"
	"github.com/pion/webrtc/v2"
	"sync"
	"time"
)

const (
	rtcpPLIInterval = time.Second * 3
)

type WebRTCMeting struct {
	uuidGenerator *meet.UUIDGenerator

	participants []*WebRTCParticipant
	mutex        sync.RWMutex

	log logging.Logger
}

func NewWebRTCMeting() *WebRTCMeting {
	return &WebRTCMeting{
		uuidGenerator: meet.NewUUIDGenerator(),
		log:           logging.New("meeting"),
	}
}

func (m *WebRTCMeting) Join(participant *domain.Participant) (*WebRTCParticipant, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//offer, err := decodeSdp(sdp)
	//if err != nil {
	//	return nil, errors.Wrap(err, "decoding session description protocol error")
	//}

	webrtcParticipant, err := NewWebRTCParticipant(participant)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a WebRTC participant")
	}

	m.participants = append(m.participants, webrtcParticipant)

	go func() {
		for {
			track, ok := <-webrtcParticipant.Tracks()
			if !ok {
				return
			}
			m.addTrackToMembers(track)
		}
	}()

	return webrtcParticipant, nil
}

func (m *WebRTCMeting) addTrackToMembers(track *webrtc.Track) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	m.log.Debug("adding a track to participants")

	for _, member := range m.participants {
		if err := member.AddTrack(track); err != nil {
			m.log.Warn("could not add a track", "err", err)
		}
	}
}
