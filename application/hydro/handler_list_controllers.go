package hydro

import (
	"context"

	"github.com/boreq/hydro/domain"
)

type ListControllersHandler struct {
	transactionProvider TransactionProvider
}

func NewListControllersHandler(transactionProvider TransactionProvider) *ListControllersHandler {
	return &ListControllersHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *ListControllersHandler) Execute(ctx context.Context) (controllers []*domain.Controller, err error) {
	err = h.transactionProvider.Transact(ctx, func(t *TransactableAdapters) error {
		controllers, err = t.Controllers.List()
		return err
	})
	return controllers, err
}
