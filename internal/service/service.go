package service

import (
	"context"
	"net/http"

	"github.com/boreq/errors"
	httpPort "github.com/boreq/meet/ports/http"
)

type Service struct {
	httpServer *httpPort.Server

	ctx    context.Context
	cancel context.CancelFunc

	errC           chan error
	startedCounter int
}

func NewService(httpServer *httpPort.Server) *Service {
	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		httpServer: httpServer,
		ctx:        ctx,
		cancel:     cancel,
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

	//s.startedCounter++
	//go func() {
	//	if err := s.scanner.Run(s.ctx); err != nil {
	//		s.errC <- errors.Wrap(err, "scanner error")
	//	} else {
	//		s.errC <- nil
	//	}
	//}()

	return nil
}

func (s *Service) Close() error {
	s.cancel()
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
