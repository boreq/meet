package meet

import (
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/boreq/meet/internal/logging"

	"github.com/boreq/errors"
	"github.com/boreq/meet/domain"
	"github.com/pion/webrtc/v2"
)

const (
	rtcpPLIInterval = time.Second * 3
)

type WebRTCMeting struct {
	uuidGenerator *UUIDGenerator

	members []*Member
	mutex   sync.RWMutex

	log logging.Logger
}

func NewWebRTCMeting() *WebRTCMeting {
	return &WebRTCMeting{
		uuidGenerator: NewUUIDGenerator(),
		log:           logging.New("meeting"),
	}
}

func (m *WebRTCMeting) Join(sdp string) (*Member, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	offer, err := decodeSdp(sdp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding session description protocol error")
	}

	member, err := m.newMember(offer)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a member")
	}

	m.members = append(m.members, member)

	go func() {
		for {
			track, ok := <-member.Tracks()
			if !ok {
				return
			}
			m.addTrackToMembers(track)
		}
	}()

	return member, nil
}

func (m *WebRTCMeting) newMember(offer webrtc.SessionDescription) (*Member, error) {
	uuid, err := m.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "uuid generation failed")
	}

	participantUUID, err := domain.NewParticipantUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a participant uuid")
	}

	return NewMember(participantUUID, offer)
}

func (m *WebRTCMeting) addTrackToMembers(track *webrtc.Track) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	m.log.Debug("adding a track to members")

	for _, member := range m.members {
		if err := member.AddTrack(track); err != nil {
			m.log.Warn("could not add a track", "err", err)
		}
	}
}

func encodeSdp(obj webrtc.SessionDescription) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", errors.Wrap(err, "json marshal failed")
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func decodeSdp(sdp string) (webrtc.SessionDescription, error) {
	var offer webrtc.SessionDescription

	b, err := base64.StdEncoding.DecodeString(sdp)
	if err != nil {
		return offer, errors.Wrap(err, "base64 decoding failed")
	}

	if err = json.Unmarshal(b, &offer); err != nil {
		return offer, errors.Wrap(err, "json unmarshal failed")
	}

	return offer, nil
}
