package domain

import (
	"github.com/boreq/errors"
)

type ParticipantName struct {
	name string
}

func NewParticipantName(name string) (ParticipantName, error) {
	if len(name) > 100 {
		return ParticipantName{}, errors.New("participant name too long")
	}

	return ParticipantName{name}, nil
}

func MustNewParticipantName(name string) ParticipantName {
	v, err := NewParticipantName(name)
	if err != nil {
		panic(err)
	}
	return v
}

func (n ParticipantName) String() string {
	return n.name
}
