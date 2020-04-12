package domain

import (
	"sync"

	"github.com/boreq/errors"
)

type Meeting struct {
	participants []*Participant
	mutex        sync.RWMutex

	broadcastC chan Message
	closeC     chan struct{}
}

func NewMeeting() (*Meeting, error) {
	m := &Meeting{
		broadcastC: make(chan Message),
	}

	go m.run()

	return m, nil
}

func (m *Meeting) Join(uuid ParticipantUUID, send chan<- Message) (*Participant, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	participant, err := NewParticipant(uuid, send, m.broadcastC)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a participant")
	}

	for _, existingParticipant := range m.participants {
		existingParticipant.sync(participant)
	}

	m.participants = append(m.participants, participant)

	go func() {
		<-participant.Closed()
		m.remove(participant)
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

func (m *Meeting) broadcast(msg Message) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, participant := range m.participants {
		participant.send(msg)
	}
}

func (m *Meeting) remove(p *Participant) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := range m.participants {
		if m.participants[i].UUID() == p.UUID() {
			m.participants = append(m.participants[:i], m.participants[i+1:]...)
			return
		}
	}
}
