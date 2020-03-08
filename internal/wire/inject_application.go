package wire

import (
	"github.com/boreq/hydro/application"
	"github.com/boreq/hydro/application/auth"
	"github.com/boreq/hydro/application/hydro"
	"github.com/google/wire"
)

//lint:ignore U1000 because
var appSet = wire.NewSet(
	wire.Struct(new(application.Application), "*"),

	wire.Struct(new(auth.Auth), "*"),
	auth.NewRegisterInitialHandler,
	auth.NewLoginHandler,
	auth.NewLogoutHandler,
	auth.NewCheckAccessTokenHandler,
	auth.NewListHandler,
	auth.NewCreateInvitationHandler,
	auth.NewRegisterHandler,
	auth.NewRemoveHandler,
	auth.NewSetPasswordHandler,

	wire.Struct(new(hydro.Hydro), "*"),
	hydro.NewAddControllerHandler,
	hydro.NewListControllersHandler,
)
