package controllers

import (
	"encoding/json"
	"goshare/internal/config"
	"goshare/internal/resources"
	auth "goshare/pkg/ecdsa"
	"goshare/pkg/otp"
	"log"
	"net/http"
	"time"
)

type sendResponse struct {
	Token          string
	GeneratingTime time.Time
	Duration       time.Duration
}

func (c *Controller) BroadcastHandler(w http.ResponseWriter, r *http.Request) {
	req := &resources.BroadcastRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request"))
		return
	}

	if len(req.CurrentDeviceId) == 0 || req.Files == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request"))
		return
	}

	if ok := auth.Verify(req.CurrentDeviceId, config.GetConfig().Key.PublicKey); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid user"))
		return
	}

	t := time.Now()
	tt := t.Add(c.PubSub.DefaultTtl)
	token, rtime, err := otp.GenerateOTPCode("token", tt, 6)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("err generating token", err)
		return
	}

	log.Printf("token %s generated: %d time left\n", token, rtime)
	c.PubSub.Set(token)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sendResponse{
		Token:          token,
		GeneratingTime: t,
		Duration:       c.PubSub.DefaultTtl,
	}); err != nil {
		log.Println("error marshling send response", err)
	}
}
