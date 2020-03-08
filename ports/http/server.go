package http

import (
	"net"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/boreq/errors"
	"github.com/boreq/hydro/internal/logging"
	"github.com/rs/cors"
)

type Server struct {
	server   *http.Server
	listener net.Listener
	log      logging.Logger
}

func NewServer(handler http.Handler, address string) (*Server, error) {
	log := logging.New("ports/http.Server")

	log.Info("starting listening", "address", address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, errors.Wrapf(err, "could not listen on '%s'", address)
	}

	return &Server{
		server: &http.Server{
			Handler: addMiddlewares(handler),
		},
		listener: listener,
		log:      log,
	}, nil
}

func (s *Server) Serve() error {
	s.log.Debug("starting processing connections")
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
