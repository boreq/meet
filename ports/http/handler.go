package http

import (
	"github.com/boreq/hydro/application"
	"github.com/boreq/hydro/internal/logging"
	"github.com/boreq/hydro/ports/http/frontend"
	"github.com/boreq/hydro/ports/http/hydro"
	"github.com/go-chi/chi"
	"net/http"
)

type Handler struct {
	app    *application.Application
	router *chi.Mux
	log    logging.Logger
}

const hydroPrefix = "/api/hydro"

func NewHandler(app *application.Application) (*Handler, error) {
	h := &Handler{
		app:    app,
		router: chi.NewRouter(),
		log:    logging.New("ports/http.Handler"),
	}

	h.router.Mount(hydroPrefix, http.StripPrefix(hydroPrefix, hydro.NewHandler(app)))

	// Frontend
	ffs, err := frontend.NewFrontendFileSystem()
	if err != nil {
		return nil, err
	}
	h.router.NotFound(http.FileServer(ffs).ServeHTTP)

	return h, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(rw, req)
}
