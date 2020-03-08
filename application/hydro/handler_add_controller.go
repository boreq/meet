package hydro

import (
	"context"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/domain"
)

type AddController struct {
	Address domain.Address
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
	controller, err := h.newController(cmd)
	if err != nil {
		return errors.Wrap(err, "could not create a controller")
	}

	return h.transactionProvider.Transact(ctx, func(t *TransactableAdapters) error {
		_, err := t.Controllers.GetByAddress(controller.Address())
		if err != nil {
			if !errors.Is(err, ControllerNotFoundErr) {
				return errors.Wrap(err, "could not get a controller by address")
			}
		} else {
			return errors.New("controller with this address already exists")
		}

		return t.Controllers.Save(controller)
	})
}

func (h *AddControllerHandler) newController(cmd AddController) (*domain.Controller, error) {
	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate a uuid")
	}

	controllerUUID, err := domain.NewControllerUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a controller uuid")
	}

	return domain.NewController(controllerUUID, cmd.Address)
}
