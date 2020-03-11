package hydro

import (
	"github.com/boreq/hydro/application/hydro"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/eventsourcing"
)

type ControllerRepositoryMock struct {
	Events map[domain.ControllerUUID]eventsourcing.EventSourcingEvents
}

func NewControllerRepositoryMock() *ControllerRepositoryMock {
	return &ControllerRepositoryMock{
		Events: make(map[domain.ControllerUUID]eventsourcing.EventSourcingEvents),
	}
}

func (c ControllerRepositoryMock) List() ([]*domain.Controller, error) {
	panic("implement me")
}

func (c ControllerRepositoryMock) Get(uuid domain.ControllerUUID) (*domain.Controller, error) {
	events := c.Events[uuid]
	if len(events) == 0 {
		return nil, hydro.ErrControllerNotFound
	}
	return domain.NewControllerFromHistory(events)
}

func (c ControllerRepositoryMock) GetByAddress(address domain.Address) (*domain.Controller, error) {
	panic("implement me")
}

func (c ControllerRepositoryMock) Save(controller *domain.Controller) error {
	c.Events[controller.UUID()] = append(c.Events[controller.UUID()], controller.PopChanges()...)
	return nil
}
