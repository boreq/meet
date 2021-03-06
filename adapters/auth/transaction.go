package auth

import (
	"github.com/boreq/errors"
	"github.com/boreq/meet/application/auth"
	bolt "go.etcd.io/bbolt"
)

type AuthRepositoriesProvider interface {
	Provide(tx *bolt.Tx) (*auth.TransactableRepositories, error)
}

type AuthTransactionProvider struct {
	db                   *bolt.DB
	repositoriesProvider AuthRepositoriesProvider
}

func NewAuthTransactionProvider(
	db *bolt.DB,
	repositoriesProvider AuthRepositoriesProvider,
) *AuthTransactionProvider {
	return &AuthTransactionProvider{
		db:                   db,
		repositoriesProvider: repositoriesProvider,
	}
}

func (p *AuthTransactionProvider) Read(handler auth.TransactionHandler) error {
	return p.db.View(func(tx *bolt.Tx) error {
		repositories, err := p.repositoriesProvider.Provide(tx)
		if err != nil {
			return errors.Wrap(err, "could not provide the repositories")
		}
		return handler(repositories)
	})
}

func (p *AuthTransactionProvider) Write(handler auth.TransactionHandler) error {
	return p.db.Update(func(tx *bolt.Tx) error {
		repositories, err := p.repositoriesProvider.Provide(tx)
		if err != nil {
			return errors.Wrap(err, "could not provide the repositories")
		}
		return handler(repositories)
	})
}
