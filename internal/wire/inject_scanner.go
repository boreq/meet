package wire

import (
	"time"

	"github.com/boreq/hydro/application"
	"github.com/boreq/hydro/ports/scanner"
	"github.com/google/wire"
)

//lint:ignore U1000 because
var scannerSet = wire.NewSet(
	newScanner,
)

//lint:ignore U1000 because
var testScannerSet = wire.NewSet(
	newTestScanner,
)

func newScanner(client scanner.ControllerClient, app *application.Application) *scanner.Scanner {
	return scanner.NewScanner(client, app, 5*time.Minute)
}

func newTestScanner(client scanner.ControllerClient, app *application.Application) *scanner.Scanner {
	return scanner.NewScanner(client, app, 1*time.Second)
}
