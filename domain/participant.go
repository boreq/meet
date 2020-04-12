package domain

import "github.com/boreq/errors"

type Participant struct {
	uuid ParticipantUUID
	name Name

	sendC      chan<- Message
	broadcastC chan<- Message
	closeC     chan struct{}
}

func NewParticipant(uuid ParticipantUUID, sendC chan<- Message, broadcastC chan<- Message) (*Participant, error) {
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

func (p *Participant) SetName(name Name) {
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

func (p *Participant) sync(o *Participant) {
	o.send(p.nameChangedMessage())
}

func (p *Participant) send(msg Message) {
	select {
	case p.sendC <- msg:
		return
	case <-p.closeC:
		return
	}
}

func (p *Participant) broadcast(msg Message) {
	select {
	case p.broadcastC <- msg:
		return
	case <-p.closeC:
		return
	}
}

func (p *Participant) nameChangedMessage() NameChangedMessage {
	return NameChangedMessage{p.uuid, p.name}
}

type Message interface {
	// Critical returns true if this message can't be dropped in case of
	// connection problems.
	Critical() bool // todo
}

type NameChangedMessage struct {
	ParticipantUUID ParticipantUUID
	Name            Name
}

func (n NameChangedMessage) Critical() bool {
	return true
}
