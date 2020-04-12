package adapters

import (
	"github.com/boreq/errors"
	"github.com/boreq/meet/internal/eventsourcing"
	"github.com/boreq/meet/internal/eventsourcing/adapters"
	"github.com/boreq/meet/internal/eventsourcing/example/domain"
)

type BankAccountRepository struct {
	eventStore *eventsourcing.EventStore
}

func NewBankAccountRepository() *BankAccountRepository {
	persistenceAdapter := adapters.NewMemoryPersistenceAdapter()
	eventStore := eventsourcing.NewEventStore(mapping, persistenceAdapter)
	return &BankAccountRepository{
		eventStore: eventStore,
	}
}

func (r *BankAccountRepository) Save(account *domain.BankAccount) error {
	events := account.PopChanges()
	return r.eventStore.SaveEvents(r.convertUUID(account.UUID()), events)
}

func (r *BankAccountRepository) Get(uuid domain.BankAccountUUID) (*domain.BankAccount, error) {
	events, err := r.eventStore.GetEvents(r.convertUUID(uuid))
	if err != nil {
		return nil, errors.Wrap(err, "could not get the events")
	}

	return domain.NewBankAccountFromHistory(events)
}

func (r *BankAccountRepository) convertUUID(uuid domain.BankAccountUUID) eventsourcing.AggregateUUID {
	return eventsourcing.AggregateUUID(uuid.String())
}
