package domain

import "github.com/boreq/errors"

type Participant struct {
	uuid ParticipantUUID
	name ParticipantName

	sendC      chan<- OutgoingMessage
	broadcastC chan<- BroadcastMessage
	closeC     chan struct{}
}

func NewParticipant(uuid ParticipantUUID, sendC chan<- OutgoingMessage, broadcastC chan<- BroadcastMessage) (*Participant, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	return &Participant{
		uuid:       uuid,
		sendC:      sendC,
		broadcastC: broadcastC,
		closeC:     make(chan struct{}),
	}, nil
}

func (p *Participant) SetName(name ParticipantName) {
	p.name = name
	p.broadcast(p.nameChangedMessage())
}

func (p *Participant) UUID() ParticipantUUID {
	return p.uuid
}

func (p *Participant) Closed() <-chan struct{} {
	return p.closeC
}

func (p *Participant) Close() error {
	close(p.closeC)
	return nil
}

func (p *Participant) syncJoin(newParticipant *Participant) {
	// joined
	newParticipant.send(p.joinedMessage())
	p.send(newParticipant.joinedMessage())

	// name changed
	newParticipant.send(p.nameChangedMessage())
}

func (p *Participant) syncQuit(quittingParticipant *Participant) {
	// quit
	p.send(quittingParticipant.quitMessage())
}

func (p *Participant) send(msg OutgoingMessage) {
	select {
	case p.sendC <- msg:
		return
	case <-p.closeC:
		return
	}
}

func (p *Participant) broadcast(msg OutgoingMessage) {
	broadcastMessage := BroadcastMessage{
		Sender:  p.uuid,
		Message: msg,
	}

	select {
	case p.broadcastC <- broadcastMessage:
		return
	case <-p.closeC:
		return
	}
}

func (p *Participant) joinedMessage() JoinedMessage {
	return JoinedMessage{p.uuid}
}

func (p *Participant) nameChangedMessage() NameChangedMessage {
	return NameChangedMessage{p.uuid, p.name}
}

func (p *Participant) quitMessage() QuitMessage {
	return QuitMessage{p.uuid}
}
