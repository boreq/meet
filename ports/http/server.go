package http

import (
	"github.com/boreq/errors"
	"net"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/boreq/hydro/internal/logging"
	"github.com/rs/cors"
)

type Server struct {
	server   *http.Server
	listener net.Listener
	log      logging.Logger
}

func NewServer(handler http.Handler, address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, errors.Wrapf(err, "could not listen on '%s'", address)
	}

	return &Server{
		server: &http.Server{
			Handler: addMiddlewares(handler),
		},
		listener: listener,
		log:      logging.New("ports/http.Server"),
	}, nil
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Close() error {
	return s.server.Close()
}

func addMiddlewares(handler http.Handler) http.Handler {
	// Add CORS middleware
	handler = cors.AllowAll().Handler(handler)

	// Add GZIP middleware
	handler = gziphandler.GzipHandler(handler)

	return handler
}
