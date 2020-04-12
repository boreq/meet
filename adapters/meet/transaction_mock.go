package meet

import (
	"context"

	"github.com/boreq/meet/application/meet"
)

type MockTransactionProvider struct {
	transactableAdapters *meet.TransactableAdapters
}

func NewMockTransactionProvider(transactableAdapters *meet.TransactableAdapters) *MockTransactionProvider {
	return &MockTransactionProvider{
		transactableAdapters: transactableAdapters,
	}
}

func (p *MockTransactionProvider) Transact(ctx context.Context, handler meet.TransactionHandler) error {
	return handler(p.transactableAdapters)
}
