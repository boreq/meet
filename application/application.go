package application

import (
	"github.com/boreq/hydro/application/auth"
	"github.com/boreq/hydro/application/hydro"
)

type Application struct {
	Auth auth.Auth
	Hydro hydro.Hydro
}
