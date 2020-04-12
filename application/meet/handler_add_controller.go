package meet

import (
	"context"
)

type AddController struct {
}

type AddControllerHandler struct {
	transactionProvider TransactionProvider
	uuidGenerator       UUIDGenerator
}

func NewAddControllerHandler(
	transactionProvider TransactionProvider,
	uuidGenerator UUIDGenerator,
) *AddControllerHandler {
	return &AddControllerHandler{
		transactionProvider: transactionProvider,
		uuidGenerator:       uuidGenerator,
	}
}

func (h *AddControllerHandler) Execute(ctx context.Context, cmd AddController) error {
	return nil
}
