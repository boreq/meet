package meet

import (
	"context"

	"github.com/boreq/errors"
	"github.com/boreq/meet/application/meet"
	bolt "go.etcd.io/bbolt"
)

type AdaptersProvider interface {
	Provide(tx *bolt.Tx) (*meet.TransactableAdapters, error)
}

type TransactionProvider struct {
	db               *bolt.DB
	adaptersProvider AdaptersProvider
}

func NewTransactionProvider(db *bolt.DB, adaptersProvider AdaptersProvider) *TransactionProvider {
	return &TransactionProvider{
		db:               db,
		adaptersProvider: adaptersProvider,
	}
}

func (p *TransactionProvider) Transact(ctx context.Context, handler meet.TransactionHandler) error {
	return p.db.Update(func(tx *bolt.Tx) error {
		repositories, err := p.adaptersProvider.Provide(tx)
		if err != nil {
			return errors.Wrap(err, "could not provide the repositories")
		}

		return handler(repositories)
	})
}
