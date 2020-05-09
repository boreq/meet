package domain

import (
	"sync"

	"github.com/boreq/errors"
)

type BroadcastMessage struct {
	Sender  ParticipantUUID
	Message OutgoingMessage
}

type Meeting struct {
	participants []*Participant
	mutex        sync.RWMutex

	broadcastC chan BroadcastMessage
	closeC     chan struct{}
}

func NewMeeting() *Meeting {
	m := &Meeting{
		broadcastC: make(chan BroadcastMessage),
	}

	go m.run()

	return m
}

func (m *Meeting) Join(uuid ParticipantUUID, send chan<- OutgoingMessage) (*Participant, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	participant, err := NewParticipant(uuid, send, m.broadcastC)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a participant")
	}

	m.syncJoin(participant)

	m.participants = append(m.participants, participant)

	go func() {
		<-participant.Closed()
		m.onParticipantClosed(participant)
	}()

	return participant, nil
}

func (m *Meeting) Close() error {
	close(m.closeC)
	return nil
}

func (m *Meeting) run() {
	for {
		select {
		case msg := <-m.broadcastC:
			m.broadcast(msg)
		case <-m.closeC:
			return
		}
	}
}

func (m *Meeting) broadcast(msg BroadcastMessage) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, participant := range m.participants {
		if participant.UUID() != msg.Sender {
			participant.send(msg.Message)
		}
	}
}

func (m *Meeting) onParticipantClosed(participant *Participant) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.remove(participant)
	m.syncQuit(participant)
}

func (m *Meeting) remove(p *Participant) {
	for i := range m.participants {
		if m.participants[i].UUID() == p.UUID() {
			m.participants = append(m.participants[:i], m.participants[i+1:]...)
			return
		}
	}
}

func (m *Meeting) syncQuit(participant *Participant) {
	for _, remainingParticipant := range m.participants {
		remainingParticipant.syncQuit(participant)
	}
}

func (m *Meeting) syncJoin(participant *Participant) {
	for _, existingParticipant := range m.participants {
		if participant.uuid != existingParticipant.uuid {
			existingParticipant.syncJoin(participant)
		}
	}
}

func (m *Meeting) PassSessionDescription(targetParticipantUUID ParticipantUUID, sourceParticipantUUID ParticipantUUID, sessionDescription SessionDescription) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	msg := RemoteSessionDescription{
		ParticipantUUID:    sourceParticipantUUID,
		SessionDescription: sessionDescription,
	}

	m.sendTo(targetParticipantUUID, msg)
}

func (m *Meeting) PassIceCandidate(targetParticipantUUID ParticipantUUID, sourceParticipantUUID ParticipantUUID, iceCandidate IceCandidate) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	msg := RemoteIceCandidate{
		ParticipantUUID: sourceParticipantUUID,
		IceCandidate:    iceCandidate,
	}

	m.sendTo(targetParticipantUUID, msg)
}

func (m *Meeting) sendTo(participantUUID ParticipantUUID, msg OutgoingMessage) {
	for _, existingParticipant := range m.participants {
		if participantUUID == existingParticipant.uuid {
			existingParticipant.send(msg)
		}
	}
}
