//+build wireinject

package wire

import (
	"github.com/boreq/hydro/application/auth"
	"github.com/boreq/hydro/application/hydro"
	"github.com/boreq/hydro/internal/config"
	"github.com/boreq/hydro/internal/service"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

func BuildTransactableAuthRepositories(tx *bolt.Tx) (*auth.TransactableRepositories, error) {
	wire.Build(
		adapterSet,
	)

	return nil, nil
}

func BuildTransactableHydroAdapters(tx *bolt.Tx) (*hydro.TransactableAdapters, error) {
	wire.Build(
		adapterSet,
	)

	return nil, nil
}

func BuildAuthForTest(db *bolt.DB) (*auth.Auth, error) {
	wire.Build(
		appSet,
		adapterSet,
	)

	return nil, nil
}

func BuildAuth(conf *config.Config) (*auth.Auth, error) {
	wire.Build(
		appSet,
		boltSet,
		adapterSet,
	)

	return nil, nil
}

func BuildService(conf *config.Config) (*service.Service, error) {
	wire.Build(
		service.NewService,
		httpSet,
		appSet,
		boltSet,
		adapterSet,
	)

	return nil, nil
}

func BuildComponentTestService(db *bolt.DB, conf *config.Config) (ComponentTestService, error) {
	wire.Build(
		service.NewService,
		httpSet,
		appSet,
		adapterSet,

		wire.Struct(new(ComponentTestService), "*"),
	)

	return ComponentTestService{}, nil
}

type ComponentTestService struct {
	Service *service.Service
	Config  *config.Config
}
