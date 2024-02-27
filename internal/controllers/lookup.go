package controllers

import (
	"encoding/json"
	"goshare/internal/resources"
	"log"
	"net/http"
	"strings"
)

func (c Controller) LookupHandler(w http.ResponseWriter, r *http.Request) {
	req := resources.LookupRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request"))
		return
	}

	if _, err := extractDeviceToken(r, c.crypto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if strings.TrimSpace(req.Otp) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid otp"))
		return
	}

	data, ok := c.pubSub.Get(req.Otp)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("resource not found"))
		return
	}

	files, ok := data.([]resources.File)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.Println("error marshling send response", err)
		return
	}
}
