package hydro

import (
	"encoding/json"
	"net/http"

	"github.com/boreq/hydro/application"
	"github.com/boreq/hydro/application/hydro"
	"github.com/boreq/hydro/domain"
	"github.com/boreq/hydro/internal/logging"
	"github.com/boreq/rest"
	"github.com/go-chi/chi"
)

type Handler struct {
	app    *application.Application
	router *chi.Mux
	log    logging.Logger
}

func NewHandler(app *application.Application) *Handler {
	h := &Handler{
		app:    app,
		router: chi.NewRouter(),
		log:    logging.New("ports/http/hydro"),
	}

	h.registerHandlers()

	return h
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(rw, req)
}

func (h *Handler) registerHandlers() {
	h.router.Get("/controllers", rest.Wrap(h.listControllers))
	h.router.Get("/controllers/{controllerUUID}/devices", rest.Wrap(h.listControllerDevices))
	h.router.Post("/controllers", rest.Wrap(h.addController))
}

func (h *Handler) listControllers(r *http.Request) rest.RestResponse {
	controllers, err := h.app.Hydro.ListControllersHandler.Execute(r.Context())
	if err != nil {
		h.log.Warn("list controllers failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(toControllers(controllers))
}

func (h *Handler) listControllerDevices(r *http.Request) rest.RestResponse {
	stringControllerUUID := chi.URLParam(r, "controllerUUID")

	controllerUUID, err := domain.NewControllerUUID(stringControllerUUID)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid controller UUID.")
	}

	query := hydro.ListControllerDevices{
		ControllerUUID: controllerUUID,
	}

	devices, err := h.app.Hydro.ListControllerDevicesHandler.Execute(r.Context(), query)
	if err != nil {
		h.log.Warn("list controller devices failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(toDevices(devices))

}

func (h *Handler) addController(r *http.Request) rest.RestResponse {
	var input AddControllerRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return rest.ErrBadRequest.WithMessage("Malformed input.")
	}

	address, err := domain.NewAddress(input.Address)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid address.")
	}

	cmd := hydro.AddController{
		Address: address,
	}

	if err := h.app.Hydro.AddControllerHandler.Execute(r.Context(), cmd); err != nil {
		h.log.Warn("add controllers failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil).WithStatusCode(http.StatusCreated)
}
