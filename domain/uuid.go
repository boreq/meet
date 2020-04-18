package domain

import "github.com/boreq/errors"

type uuid struct {
	u string
}

func newUUID(u string) (uuid, error) {
	if u == "" {
		return uuid{}, errors.New("uuid can not be empty")
	}

	return uuid{u: u}, nil
}

func (u uuid) IsZero() bool {
	return u.u == ""
}

func (u uuid) String() string {
	return u.u
}

type ParticipantUUID struct {
	uuid
}

func NewParticipantUUID(u string) (ParticipantUUID, error) {
	uuid, err := newUUID(u)
	if err != nil {
		return ParticipantUUID{}, errors.New("could not create a participant UUID")
	}
	return ParticipantUUID{uuid}, nil
}

func MustNewParticipantUUID(u string) ParticipantUUID {
	v, err := NewParticipantUUID(u)
	if err != nil {
		panic(err)
	}
	return v
}
