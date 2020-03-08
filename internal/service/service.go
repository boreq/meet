package service

import (
	"net/http"

	"github.com/boreq/errors"
	httpPort "github.com/boreq/hydro/ports/http"
)

type Service struct {
	httpServer *httpPort.Server

	errC           chan error
	startedCounter int
}

func NewService(httpServer *httpPort.Server) *Service {
	return &Service{
		httpServer: httpServer,
	}
}

func (s *Service) Start() error {
	if s.errC != nil {
		return errors.New("already started")
	}

	s.errC = make(chan error)

	s.startedCounter++
	go func() {
		if err := s.httpServer.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errC <- errors.Wrap(err, "http server error")
		} else {
			s.errC <- nil
		}
	}()

	return nil
}

func (s *Service) Close() error {
	return s.httpServer.Close()
}

func (s *Service) Wait() error {
	for i := 0; i < s.startedCounter; i++ {
		if err := <-s.errC; err != nil {
			return err
		}
	}

	return nil
}
