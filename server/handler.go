package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"

	"goshare/resources"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type getConstantResponse struct {
	P []byte
	G []byte
}

func (s *Server) handleConstants(w http.ResponseWriter, r *http.Request) error {
	err := json.NewEncoder(w).Encode(getConstantResponse{
		P: s.config.Crypto.P,
		G: s.config.Crypto.G,
	})
	if err != nil {
		return err
	}

	return nil
}

type channelRequest struct {
	Code  string            `json:"code"`
	Files []resources.Media `json:"files"`
}

func (s *Server) handleChannel(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	var request *channelRequest
	if err := conn.ReadJSON(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return errors.Wrap(err, "json.Decode")
	}

	events := s.getSetter.NotifyEvent(request.Code, request.Files)
	for event := range events {
		if err := conn.WriteJSON(event); err != nil {
			slog.Warn("error while encoding event", slog.Any("event", event))
			continue
		}
	}

	return nil
}

func (s *Server) handleSender(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("otp")
	if id == "" {
		http.Error(w, "no id passed", http.StatusBadRequest)
		return nil
	}

	mediatype, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		slog.Debug("Context-Type not set header")
		return nil
	}

	if mediatype != "multipart/form-data" {
		slog.Debug("mediatype not multipart/form-data")
		return nil
	}

	boundary, ok := params["boundary"]
	if !ok {
		slog.Debug("boundary not set")
		return nil
	}

	reader := multipart.NewReader(r.Body, boundary)

	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return errors.Wrap(err, "reader.nextPart")
		}

		defer part.Close()
		s.getSetter.SetData(id, part)
	}

	return nil
}

type getMetaDataResponse struct {
	MetaData []resources.Media
	Code     string
}

func (s *Server) handleMetaData(w http.ResponseWriter, r *http.Request) error {
	otp := r.PathValue("otp")
	if otp == "" {
		http.Error(w, "no otp passed", http.StatusBadRequest)
		return nil
	}

	metadata, code := s.getSetter.GetMetaData(otp)

	err := json.NewEncoder(w).Encode(getMetaDataResponse{
		MetaData: metadata,
		Code:     code,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleRecvMedia(w http.ResponseWriter, r *http.Request) error {
	otp := r.PathValue("otp")
	if otp == "" {
		http.Error(w, "otp cannot be empty", http.StatusBadRequest)
		return nil
	}

	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "code cannot be empty", http.StatusBadRequest)
		return nil
	}

	reader, err := s.getSetter.GetData(otp, code)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return err
	}

	_, err = io.Copy(w, reader)
	return err
}
