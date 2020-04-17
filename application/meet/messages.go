package meet

import "github.com/boreq/meet/domain"

type IncomingMessage interface{}

type SetNameMessage struct {
	Name domain.Name
}
