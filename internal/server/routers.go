package server

import (
	"goshare/internal/controllers"
	"net/http"
)

type Server struct {
	controller *controllers.Controller
}

func NewServer(controller *controllers.Controller) *Server {
	s := &Server{
		controller: controller,
	}
	s.setRoutes()
	return s
}

func (s *Server) setRoutes() {
	http.HandleFunc("POST /broadcast", s.controller.BroadcastHandler)
	http.HandleFunc("GET /identity", s.controller.IdentityHandler)
	http.HandleFunc("GET /lookup", s.controller.LookupHandler)
}

func (s *Server) Start(addr string) error {
	if err := http.ListenAndServe(addr, nil); err != http.ErrServerClosed {
		return err
	}

	return nil
}
