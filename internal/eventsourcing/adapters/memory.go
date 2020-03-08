package adapters

import "github.com/boreq/hydro/internal/eventsourcing"

type MemoryPersistenceAdapter struct {
	events map[eventsourcing.AggregateUUID][]eventsourcing.PersistedEvent
}

func NewMemoryPersistenceAdapter() *MemoryPersistenceAdapter {
	return &MemoryPersistenceAdapter{
		events: make(map[eventsourcing.AggregateUUID][]eventsourcing.PersistedEvent),
	}
}

func (m *MemoryPersistenceAdapter) SaveEvents(uuid eventsourcing.AggregateUUID, events []eventsourcing.PersistedEvent) error {
	m.events[uuid] = append(m.events[uuid], events...)
	return nil
}

func (m *MemoryPersistenceAdapter) GetEvents(uuid eventsourcing.AggregateUUID) ([]eventsourcing.PersistedEvent, error) {
	return m.events[uuid], nil
}
