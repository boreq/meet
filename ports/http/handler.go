package http

import (
	"net/http"

	"github.com/boreq/meet/application"
	"github.com/boreq/meet/ports/http/frontend"
	"github.com/boreq/meet/ports/http/meet"
	"github.com/go-chi/chi"
)

const meetPrefix = "/api/meet"

type Handler struct {
	router *chi.Mux
}

func NewHandler(app *application.Application) (*Handler, error) {
	h := &Handler{
		router: chi.NewRouter(),
	}

	// Subrouters
	h.router.Mount(meetPrefix, http.StripPrefix(meetPrefix, meet.NewHandler(app)))

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
