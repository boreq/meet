package eventsourcing

type EventSourcing struct {
	Changes        []EventSourcingEvent
	CurrentVersion AggregateVersion
}

func (e *EventSourcing) Record(event Event) error {
	e.Changes = append(e.Changes, EventSourcingEvent{
		Event:            event,
		AggregateVersion: e.CurrentVersion,
	})
	e.CurrentVersion += 1
	return nil
}

func (e *EventSourcing) HasChanges() bool {
	return len(e.Changes) > 0
}

func (e *EventSourcing) PopChanges() []EventSourcingEvent {
	tmp := e.Changes
	e.Changes = nil
	return tmp
}

func (e *EventSourcing) LoadVersion(event EventSourcingEvent) {
	e.CurrentVersion = event.AggregateVersion + 1
}
