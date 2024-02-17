package http

import (
	"goshare/internal/config"
	"goshare/internal/controllers"
	"log"
	"net/http"
)

type Server struct {
	controller *controllers.Controller
}

func NewServer(controller *controllers.Controller) *Server {
	s := &Server{
		controller: controller,
	}
	s.Init()
	return s
}

func (s *Server) Init() {
	http.HandleFunc("/send", s.controller.BroadcastHandler)
}

func (s *Server) start() {
	cfg := config.GetConfig()
	if err := http.ListenAndServe(cfg.Server.Addr, nil); err != http.ErrServerClosed {
		log.Println(err)
	}
}
