package hydro

import (
	"context"

	"github.com/boreq/errors"

	"github.com/boreq/hydro/domain"
)

type ListControllerDevices struct {
	ControllerUUID domain.ControllerUUID
}

type ListControllerDevicesHandler struct {
	transactionProvider TransactionProvider
}

func NewListControllerDevicesHandler(transactionProvider TransactionProvider) *ListControllerDevicesHandler {
	return &ListControllerDevicesHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *ListControllerDevicesHandler) Execute(ctx context.Context, query ListControllerDevices) (devices []*domain.Device, err error) {
	err = h.transactionProvider.Transact(ctx, func(t *TransactableAdapters) error {
		_, err = t.Controllers.Get(query.ControllerUUID) // todo remove?
		if err != nil {
			if errors.Is(err, ErrControllerNotFound) {
				return err
			}
			return errors.Wrap(err, "error retrieving the controller")
		}

		devices, err = t.Devices.ListByController(query.ControllerUUID)
		return err
	})
	return devices, err
}
