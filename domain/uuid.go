package domain

import "github.com/boreq/errors"

type uuid struct {
	u string
}

func (u uuid) IsZero() bool {
	return u.u == ""
}

func newUUID(u string) (uuid, error) {
	if u == "" {
		return uuid{}, errors.New("uuid can not be empty")
	}

	return uuid{u: u}, nil
}

type RackUUID struct {
	uuid
}

func NewRackUUID(u string) (RackUUID, error) {
	uuid, err := newUUID(u)
	if err != nil {
		return RackUUID{}, errors.New("could not create a rack UUID")
	}
	return RackUUID{uuid}, nil
}

type DeviceUUID struct {
	uuid
}

func NewDeviceUUID(u string) (DeviceUUID, error) {
	uuid, err := newUUID(u)
	if err != nil {
		return DeviceUUID{}, errors.New("could not create a device UUID")
	}
	return DeviceUUID{uuid}, nil
}
