package hydro

import "github.com/boreq/hydro/domain"

type Controller struct {
	UUID    string `json:"uuid"`
	Address string `json:"address"`
}

func toControllers(controllers []*domain.Controller) []Controller {
	rv := make([]Controller, 0)
	for _, controller := range controllers {
		rv = append(rv, toController(controller))
	}
	return rv
}

func toController(controller *domain.Controller) Controller {
	return Controller{
		UUID:    controller.UUID().String(),
		Address: controller.Address().String(),
	}
}
