package domain

import "github.com/boreq/errors"

type SessionDescription struct {
	d string
}

func NewSessionDescription(d string) (SessionDescription, error) {
	if d == "" {
		return SessionDescription{}, errors.New("description can not be empty")
	}

	return SessionDescription{d: d}, nil
}

func (d SessionDescription) IsZero() bool {
	return d.d == ""
}

func (d SessionDescription) String() string {
	return d.d
}
