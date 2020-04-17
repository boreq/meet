package domain

import (
	"github.com/boreq/errors"
)

type Name struct {
	name string
}

func NewName(name string) (Name, error) {
	if len(name) > 100 {
		return Name{}, errors.New("name too long")
	}

	return Name{name}, nil
}

func MustNewName(name string) Name {
	v, err := NewName(name)
	if err != nil {
		panic(err)
	}
	return v
}

func (n Name) String() string {
	return n.name
}
