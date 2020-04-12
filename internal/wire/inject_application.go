package wire

import (
	"github.com/boreq/meet/application"
	"github.com/boreq/meet/application/auth"
	"github.com/boreq/meet/application/meet"
	"github.com/google/wire"
)

//lint:ignore U1000 because
var appSet = wire.NewSet(
	wire.Struct(new(application.Application), "*"),
	authAppSet,
	meetAppSet,
)

var authAppSet = wire.NewSet(
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
)

var meetAppSet = wire.NewSet(
	wire.Struct(new(meet.Meet), "*"),
)
