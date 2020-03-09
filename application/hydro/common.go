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
	Get(uuid domain.ControllerUUID) (*domain.Controller, error)
	GetByAddress(address domain.Address) (*domain.Controller, error)
	Save(controller *domain.Controller) error
}

type DeviceRepository interface {
	ListByController(uuid domain.ControllerUUID) ([]*domain.Device, error)
	Remove(uuid domain.DeviceUUID) error
	Save(device *domain.Device) error
}

type TransactableAdapters struct {
	Controllers ControllerRepository
	Devices     DeviceRepository
}

type TransactionHandler func(t *TransactableAdapters) error

type TransactionProvider interface {
	Transact(context.Context, TransactionHandler) error
}

type Hydro struct {
	AddControllerHandler *AddControllerHandler
	SetControllerDevices *SetControllerDevicesHandler

	ListControllersHandler *ListControllersHandler
}
