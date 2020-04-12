package domain

import (
	"regexp"
	"strings"

	"github.com/boreq/errors"
)

var meetingNameRegexp = regexp.MustCompile("^[a-z]+$")

type MeetingName struct {
	name string
}

func NewMeetingName(name string) (MeetingName, error) {
	name = strings.ToLower(name)

	if name == "" {
		return MeetingName{}, errors.New("meeting name can not be empty")
	}

	if len(name) > 100 {
		return MeetingName{}, errors.New("meeting name too long")
	}

	if !meetingNameRegexp.MatchString(name) {
		return MeetingName{}, errors.New("invalid meeting name")
	}

	return MeetingName{name}, nil
}

func MustNewMeetingName(name string) MeetingName {
	v, err := NewMeetingName(name)
	if err != nil {
		panic(err)
	}
	return v
}

func (n MeetingName) IsZero() bool {
	return n == MeetingName{}
}
