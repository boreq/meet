package controller

import (
	"context"
	"fmt"

	"github.com/boreq/hydro/domain"
)

type ClientMock struct {
}

func NewClientMock() *ClientMock {
	return &ClientMock{}
}

func (c *ClientMock) GetState(ctx context.Context, address domain.Address) (ControllerState, error) {
	return ControllerState{
		Devices: []Device{
			mockDevice(1),
			mockDevice(2),
			mockDevice(3),
			mockDevice(4),
			mockDevice(5),
		},
	}, nil
}

func mockDevice(i int) Device {
	return Device{
		Id: fmt.Sprintf("mock-device-%d", i),
	}
}
