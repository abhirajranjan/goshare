package server

import (
	"context"
	"net/http"
)

type Server struct {
	server    *http.Server
	err       error
	isClosing bool
	Wait      chan struct{}
}

func (s *Server) Start() {
	go func() {
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			s.err = err
			if !s.isClosing {
				s.Wait <- struct{}{}
			}
		}
	}()
}

func (s *Server) Close() {
	go func() {
		s.server.RegisterOnShutdown(func() { s.isClosing = true })
		err := s.server.Shutdown(context.Background())
		if err != nil {
			s.err = err
		}
		s.Wait <- struct{}{}
	}()
}

func (s *Server) Error() error {
	return s.err
}

func NewServer(srv *http.Server) *Server {
	return &Server{
		server: srv,
		Wait:   make(chan struct{}),
	}
}
