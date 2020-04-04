package domain

import (
	"fmt"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/internal/eventsourcing"
)

type Controller struct {
	uuid    ControllerUUID
	address Address

	es eventsourcing.EventSourcing
}

func NewController(uuid ControllerUUID, address Address) (*Controller, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	if address.IsZero() {
		return nil, errors.New("zero value of address")
	}

	controller := &Controller{}

	event := ControllerCreated{uuid, address}
	if err := controller.update(event); err != nil {
		return nil, errors.Wrap(err, "could not consume an event")
	}

	return controller, nil
}

func NewControllerFromHistory(events []eventsourcing.EventSourcingEvent) (*Controller, error) {
	controller := &Controller{}

	for _, event := range events {
		if err := controller.update(event.Event); err != nil {
			return nil, errors.Wrap(err, "could not process an event")
		}
		controller.es.LoadVersion(event)
	}

	controller.es.PopChanges()

	return controller, nil
}

func (c *Controller) UUID() ControllerUUID {
	return c.uuid
}

func (c *Controller) Address() Address {
	return c.address
}

func (c *Controller) HasChanges() bool {
	return c.es.HasChanges()
}

func (c *Controller) PopChanges() eventsourcing.EventSourcingEvents {
	return c.es.PopChanges()
}

func (c *Controller) update(event eventsourcing.Event) error {
	switch e := event.(type) {
	case ControllerCreated:
		c.handleCreated(e)
	default:
		return fmt.Errorf("unsupported event '%T'", event)
	}

	return c.es.Record(event)
}

func (c *Controller) handleCreated(e ControllerCreated) {
	c.uuid = e.UUID
	c.address = e.Address
}

type ControllerCreated struct {
	UUID    ControllerUUID
	Address Address
}

func (c ControllerCreated) EventType() eventsourcing.EventType {
	return "created_v1"
}
