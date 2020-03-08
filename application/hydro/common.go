package hydro

import (
	"context"
	"errors"

	"github.com/boreq/hydro/domain"
)

type UUIDGenerator interface {
	Generate() (string, error)
}

var ErrControllerNotFound = errors.New("controller not found")

type ControllerRepository interface {
	List() ([]*domain.Controller, error)
	GetByAddress(address domain.Address) (*domain.Controller, error)
	Save(*domain.Controller) error
}

type TransactableAdapters struct {
	Controllers ControllerRepository
}

type TransactionHandler func(t *TransactableAdapters) error

type TransactionProvider interface {
	Transact(context.Context, TransactionHandler) error
}

type Hydro struct {
	AddControllerHandler   *AddControllerHandler
	ListControllersHandler *ListControllersHandler
}
