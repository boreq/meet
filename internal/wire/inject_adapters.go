package wire

import (
	authAdapters "github.com/boreq/hydro/adapters/auth"
	hydroAdapters "github.com/boreq/hydro/adapters/hydro"
	"github.com/boreq/hydro/application/auth"
	"github.com/boreq/hydro/application/hydro"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var adapterSet = wire.NewSet(
	// auth
	authAdapters.NewAuthTransactionProvider,
	wire.Bind(new(auth.TransactionProvider), new(*authAdapters.AuthTransactionProvider)),

	wire.Struct(new(auth.TransactableRepositories), "*"),

	newAuthRepositoriesProvider,
	wire.Bind(new(authAdapters.AuthRepositoriesProvider), new(*authRepositoriesProvider)),

	wire.Bind(new(auth.UserRepository), new(*authAdapters.UserRepository)),
	authAdapters.NewUserRepository,

	wire.Bind(new(auth.InvitationRepository), new(*authAdapters.InvitationRepository)),
	authAdapters.NewInvitationRepository,

	wire.Bind(new(auth.PasswordHasher), new(*authAdapters.BcryptPasswordHasher)),
	authAdapters.NewBcryptPasswordHasher,

	wire.Bind(new(auth.AccessTokenGenerator), new(*authAdapters.CryptoAccessTokenGenerator)),
	authAdapters.NewCryptoAccessTokenGenerator,

	authAdapters.NewCryptoStringGenerator,
	wire.Bind(new(auth.CryptoStringGenerator), new(*authAdapters.CryptoStringGenerator)),

	// hydro
	hydroAdapters.NewTransactionProvider,
	wire.Bind(new(hydro.TransactionProvider), new(*hydroAdapters.TransactionProvider)),

	wire.Struct(new(hydro.TransactableAdapters), "*"),

	newHydroAdaptersProvider,
	wire.Bind(new(hydroAdapters.AdaptersProvider), new(*hydroAdaptersProvider)),

	hydroAdapters.NewControllerRepository,
	wire.Bind(new(hydro.ControllerRepository), new(*hydroAdapters.ControllerRepository)),

	hydroAdapters.NewDeviceRepository,
	wire.Bind(new(hydro.DeviceRepository), new(*hydroAdapters.DeviceRepository)),

	hydroAdapters.NewUUIDGenerator,
	wire.Bind(new(hydro.UUIDGenerator), new(*hydroAdapters.UUIDGenerator)),
)

//lint:ignore U1000 because
var testAdapterSet = wire.NewSet(
	// hydro
	hydroAdapters.NewMockTransactionProvider,
	wire.Bind(new(hydro.TransactionProvider), new(*hydroAdapters.MockTransactionProvider)),

	wire.Struct(new(hydro.TransactableAdapters), "*"),

	hydroAdapters.NewControllerRepositoryMock,
	wire.Bind(new(hydro.ControllerRepository), new(*hydroAdapters.ControllerRepositoryMock)),

	hydroAdapters.NewDeviceRepositoryMock,
	wire.Bind(new(hydro.DeviceRepository), new(*hydroAdapters.DeviceRepositoryMock)),

	hydroAdapters.NewUUIDGeneratorMock,
	wire.Bind(new(hydro.UUIDGenerator), new(*hydroAdapters.UUIDGeneratorMock)),
)

type authRepositoriesProvider struct {
}

func newAuthRepositoriesProvider() *authRepositoriesProvider {
	return &authRepositoriesProvider{}
}

func (p *authRepositoriesProvider) Provide(tx *bolt.Tx) (*auth.TransactableRepositories, error) {
	return BuildTransactableAuthRepositories(tx)
}

type hydroAdaptersProvider struct {
}

func newHydroAdaptersProvider() *hydroAdaptersProvider {
	return &hydroAdaptersProvider{}
}

func (p *hydroAdaptersProvider) Provide(tx *bolt.Tx) (*hydro.TransactableAdapters, error) {
	return BuildTransactableHydroAdapters(tx)
}

