package domain

import "github.com/boreq/errors"

type IceCandidate struct {
	d string
}

func NewIceCandidate(d string) (IceCandidate, error) {
	if d == "" {
		return IceCandidate{}, errors.New("ice candidate can not be empty")
	}

	return IceCandidate{d: d}, nil
}

func (d IceCandidate) IsZero() bool {
	return d.d == ""
}

func (d IceCandidate) String() string {
	return d.d
}
