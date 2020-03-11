package wire

import (
	"github.com/boreq/hydro/adapters/hydro/controller"
	"github.com/boreq/hydro/ports/scanner"
	"github.com/google/wire"
)

//lint:ignore U1000 because
var clientSet = wire.NewSet(
	controller.NewClientMock,
	wire.Bind(new(scanner.ControllerClient), new(*controller.ClientMock)),
)

//lint:ignore U1000 because
var mockClientSet = wire.NewSet(
	controller.NewClientMock,
	wire.Bind(new(scanner.ControllerClient), new(*controller.ClientMock)),
)
