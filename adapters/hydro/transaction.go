package hydro

import (
	"context"
	"github.com/boreq/errors"
	"github.com/boreq/hydro/application/hydro"
	bolt "go.etcd.io/bbolt"
)

type AdaptersProvider interface {
	Provide(tx *bolt.Tx) (*hydro.TransactableAdapters, error)
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

func (p *TransactionProvider) Transact(ctx context.Context, handler hydro.TransactionHandler) error {
	return p.db.Update(func(tx *bolt.Tx) error {
		repositories, err := p.adaptersProvider.Provide(tx)
		if err != nil {
			return errors.Wrap(err, "could not provide the repositories")
		}

		return handler(repositories)
	})
}
