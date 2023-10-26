package handler

import (
	"encoding/binary"
	"encoding/json"
	"goshare/internal/models"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Options func(*Handler)

type adapter interface {
	GenerateRoom() uint64
	Broadcast(*models.BroadcastModel) error
	FetchRoom(roomID uint64) ([]byte, error)
}

type Handler struct {
	Router  *mux.Router
	Adapter adapter
}

func NewHandler(adapter adapter, opts ...Options) *Handler {
	h := &Handler{
		Router:  mux.NewRouter(),
		Adapter: adapter,
	}

	for _, op := range opts {
		op(h)
	}

	defer h.initRoutes()
	return h
}

func (h *Handler) initRoutes() {
	h.Router.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			generateRoom(h.Adapter, w, r)
		default:
			http.NotFound(w, r)
		}
	})

	h.Router.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			broadcast(h.Adapter, w, r)
		default:
			http.NotFound(w, r)
		}
	})

	h.Router.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fetchRoom(h.Adapter, w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

func generateRoom(adapter adapter, w http.ResponseWriter, r *http.Request) {
	roomID := adapter.GenerateRoom()
	w.Header().Add("Content-Type", "text/plain")

	if err := binary.Write(w, binary.BigEndian, roomID); err != nil {
		internalError(errors.Wrap(err, "binary.Write"), "", w, r)
		return
	}
}

func broadcast(adapter adapter, w http.ResponseWriter, r *http.Request) {
	var reqModel RequestBroadcastModel
	if err := json.NewDecoder(r.Body).Decode(&reqModel); err != nil {
		slog.Error("http.BadRequest", "url", r.URL, "err", err)
		badRequest("incorrect model", w, r)
		return
	}

	var model models.BroadcastModel
	model.Data = []byte(reqModel.Data)
	model.Id = reqModel.Id
	model.Proto = reqModel.Proto

	if err := adapter.Broadcast(&model); err != nil {
		if errors.Is(err, models.DomainErr{}) {
			badRequest(err.Error(), w, r)
		} else {
			internalError(errors.Wrap(err, "adapter.Broadcast"), "", w, r)
		}
		return
	}
}

func fetchRoom(adapter adapter, w http.ResponseWriter, r *http.Request) {
	var model RequestRoomModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		badRequest("incorrect model", w, r)
		return
	}

	data, err := adapter.FetchRoom(model.RoomID)
	if err != nil {
		if errors.Is(err, models.DomainErr{}) {
			badRequest(err.Error(), w, r)
		} else {
			internalError(errors.Wrap(err, "adapter.FetchRoom"), "", w, r)
		}
		return
	}

	w.Write(data)
}

func internalError(cause error, httpReason string, w http.ResponseWriter, r *http.Request) {
	slog.Error("http.Error", r.URL, cause)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(httpReason))
}

func badRequest(err string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err))
}
