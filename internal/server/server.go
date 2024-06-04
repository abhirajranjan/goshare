package server

import (
	"goshare/internal/resources"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

type GetSetter interface {
	GetMetaData(otp string) (media []resources.Media, code string)
	GetData(otp string, code string) (io.Reader, error)

	NotifyEvent(code string, media []resources.Media) <-chan resources.EventType
	SetData(id string, part *multipart.Part) error
}

type Server struct {
	config    resources.Config
	mux       *http.ServeMux
	getSetter GetSetter
}

func NewServer(cfg resources.Config, getSetter GetSetter) *Server {
	s := &Server{
		config:    cfg,
		mux:       http.NewServeMux(),
		getSetter: getSetter,
	}

	s.mux.Handle("GET /constants", errorWrapper(s.handleConstants))
	s.mux.Handle("GET /channel", errorWrapper(s.handleChannel))
	s.mux.Handle("POST /send/{id}", errorWrapper(s.handleSender))
	s.mux.Handle("GET /metadata/{otp}", errorWrapper(s.handleMetaData))
	s.mux.Handle("GET /recv/{otp}/{code}", errorWrapper(s.handleRecvMedia))

	return s
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.Server.HostPort, s.mux)
}

func errorWrapper(h func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Warn(
				err.Error(),
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
			)
		}
	})
}
