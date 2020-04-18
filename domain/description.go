package domain

import "github.com/boreq/errors"

type description struct {
	d string
}

func newDescription(d string) (description, error) {
	if d == "" {
		return description{}, errors.New("description can not be empty")
	}

	return description{d: d}, nil
}

func (d description) IsZero() bool {
	return d.d == ""
}

func (d description) String() string {
	return d.d
}

type LocalDescription struct {
	description
}

func NewLocalDescription(d string) (LocalDescription, error) {
	description, err := newDescription(d)
	if err != nil {
		return LocalDescription{}, errors.Wrap(err, "could not create a local description")
	}

	return LocalDescription{description}, nil
}

type RemoteDescription struct {
	description
}

func NewRemoteDescription(d string) (RemoteDescription, error) {
	description, err := newDescription(d)
	if err != nil {
		return RemoteDescription{}, errors.Wrap(err, "could not create a remote description")
	}

	return RemoteDescription{description}, nil
}
