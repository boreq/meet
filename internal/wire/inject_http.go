package wire

import (
	"net/http"

	"github.com/boreq/meet/internal/config"
	httpPort "github.com/boreq/meet/ports/http"
	"github.com/google/wire"
)

//lint:ignore U1000 because
var httpSet = wire.NewSet(
	newServer,

	httpPort.NewHandler,
	wire.Bind(new(http.Handler), new(*httpPort.Handler)),
)

func newServer(handler http.Handler, conf *config.Config) (*httpPort.Server, error) {
	return httpPort.NewServer(handler, conf.ServeAddress)
}
