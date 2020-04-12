package application

import (
	"github.com/boreq/meet/application/auth"
	"github.com/boreq/meet/application/meet"
)

type Application struct {
	Auth auth.Auth
	Meet meet.Meet
}
