package hydro

import (
	"context"

	"github.com/boreq/hydro/application/hydro"
)

type MockTransactionProvider struct {
	transactableAdapters *hydro.TransactableAdapters
}

func NewMockTransactionProvider(transactableAdapters *hydro.TransactableAdapters) *MockTransactionProvider {
	return &MockTransactionProvider{
		transactableAdapters: transactableAdapters,
	}
}

func (p *MockTransactionProvider) Transact(ctx context.Context, handler hydro.TransactionHandler) error {
	return handler(p.transactableAdapters)
}

