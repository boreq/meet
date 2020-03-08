package http

import (
	"net/http"

	"github.com/boreq/hydro/application"
	"github.com/boreq/hydro/ports/http/frontend"
	"github.com/boreq/hydro/ports/http/hydro"
	"github.com/go-chi/chi"
)

type Handler struct {
	router *chi.Mux
}

const hydroPrefix = "/api/hydro"

func NewHandler(app *application.Application) (*Handler, error) {
	h := &Handler{
		router: chi.NewRouter(),
	}

	// Subrouters
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
