package webrtc

import (
	"github.com/boreq/errors"
	"github.com/boreq/meet/domain"
	"github.com/pion/webrtc/v2"
)

type Track struct {
	participantUUID domain.ParticipantUUID
	track           *webrtc.Track
}

func NewTrack(participantUUID domain.ParticipantUUID, track *webrtc.Track) (Track, error) {
	if participantUUID.IsZero() {
		return Track{}, errors.New("zero value of participant UUID")
	}

	if track == nil {
		return Track{}, errors.New("nil track")
	}

	return Track{
		participantUUID: participantUUID,
		track: track,
	}, nil
}

func (t Track) IsZero() bool {
	return t == Track{}
}

func (t Track) ParticipantUUID() domain.ParticipantUUID {
	return t.participantUUID
}

func (t Track) Track() *webrtc.Track {
	return t.track
}
