package meet

import (
	"context"

	"github.com/boreq/errors"
	"github.com/boreq/meet/domain"
	"github.com/boreq/meet/internal/logging"
)

var meeting = domain.NewMeeting()

type Client struct {
	Receive <-chan IncomingMessage
	Send    chan<- domain.OutgoingMessage
}

type JoinMeeting struct {
	MeetingName domain.MeetingName
	Client      Client
}

type JoinMeetingHandler struct {
	uuidGenerator UUIDGenerator
	log           logging.Logger
}

func NewJoinMeetingHandler(uuidGenerator UUIDGenerator) *JoinMeetingHandler {
	return &JoinMeetingHandler{
		uuidGenerator: uuidGenerator,
		log:           logging.New("JoinMeetingHandler"),
	}
}

func (h *JoinMeetingHandler) Execute(ctx context.Context, cmd JoinMeeting) error {
	participant, err := h.joinMeeting(cmd)
	if err != nil {
		return errors.Wrap(err, "could not join the meeting")
	}

	defer func() {
		if err := participant.Close(); err != nil {
			h.log.Warn("could not close the participant")
		}
	}()

	for {
		select {
		case msg, ok := <-cmd.Client.Receive:
			if !ok {
				return nil
			}
			h.dispatchMessage(meeting, participant, msg)
		case <-ctx.Done():
			return nil
		}
	}
}

func (h *JoinMeetingHandler) dispatchMessage(meeting *domain.Meeting, participant *domain.Participant, msg IncomingMessage) {
	switch m := msg.(type) {
	case SetNameMessage:
		participant.SetName(m.Name)
	case LocalSessionDescription:
		meeting.PassSessionDescription(m.TargetParticipantUUID, participant.UUID(), m.SessionDescription)
	case LocalIceCandidate:
		meeting.PassIceCandidate(m.TargetParticipantUUID, participant.UUID(), m.IceCandidate)
	default:
		h.log.Warn("unknown message received", "msg", msg)
	}
}

func (h *JoinMeetingHandler) joinMeeting(cmd JoinMeeting) (*domain.Participant, error) {
	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate a uuid")
	}

	participantUUID, err := domain.NewParticipantUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a participant uuid")
	}

	return meeting.Join(participantUUID, cmd.Client.Send)
}
