package wire

import (
	authAdapters "github.com/boreq/meet/adapters/auth"
	meetAdapters "github.com/boreq/meet/adapters/meet"
	"github.com/boreq/meet/application/auth"
	"github.com/boreq/meet/application/meet"
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
	meetAdapters.NewTransactionProvider,
	wire.Bind(new(meet.TransactionProvider), new(*meetAdapters.TransactionProvider)),

	wire.Struct(new(meet.TransactableAdapters), "*"),

	//newHydroAdaptersProvider,
	//wire.Bind(new(hydroAdapters.AdaptersProvider), new(*meetAdaptersProvider)),
	//
	//hydroAdapters.NewControllerRepository,
	//wire.Bind(new(hydro.ControllerRepository), new(*hydroAdapters.ControllerRepository)),
	//
	//hydroAdapters.NewDeviceRepository,
	//wire.Bind(new(hydro.DeviceRepository), new(*hydroAdapters.DeviceRepository)),
	//
	//hydroAdapters.NewUUIDGenerator,
	//wire.Bind(new(hydro.UUIDGenerator), new(*hydroAdapters.UUIDGenerator)),
)

//lint:ignore U1000 because
var testAdapterSet = wire.NewSet(
	// hydro
	meetAdapters.NewMockTransactionProvider,
	wire.Bind(new(meet.TransactionProvider), new(*meetAdapters.MockTransactionProvider)),

	wire.Struct(new(meet.TransactableAdapters), "*"),

	//hydroAdapters.NewControllerRepositoryMock,
	//wire.Bind(new(hydro.ControllerRepository), new(*hydroAdapters.ControllerRepositoryMock)),
	//
	//hydroAdapters.NewDeviceRepositoryMock,
	//wire.Bind(new(hydro.DeviceRepository), new(*hydroAdapters.DeviceRepositoryMock)),
	//
	//hydroAdapters.NewUUIDGeneratorMock,
	//wire.Bind(new(hydro.UUIDGenerator), new(*hydroAdapters.UUIDGeneratorMock)),
)

type authRepositoriesProvider struct {
}

func newAuthRepositoriesProvider() *authRepositoriesProvider {
	return &authRepositoriesProvider{}
}

func (p *authRepositoriesProvider) Provide(tx *bolt.Tx) (*auth.TransactableRepositories, error) {
	return BuildTransactableAuthRepositories(tx)
}

type meetAdaptersProvider struct {
}

func newHydroAdaptersProvider() *meetAdaptersProvider {
	return &meetAdaptersProvider{}
}

func (p *meetAdaptersProvider) Provide(tx *bolt.Tx) (*meet.TransactableAdapters, error) {
	return BuildTransactableHydroAdapters(tx)
}
