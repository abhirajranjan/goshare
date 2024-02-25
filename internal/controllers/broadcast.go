package controllers

import (
	"encoding/json"
	"fmt"
	"goshare/internal/resources"
	"log"
	"net/http"
)

func (c Controller) BroadcastHandler(w http.ResponseWriter, r *http.Request) {
	req := resources.BroadcastRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request"))
		return
	}

	token, err := extractDeviceToken(r, c.crypto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if len(req.Files) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request"))
		return
	}

	fmt.Println(token)
	otp, originTime, interval := c.otpGen.GenerateOTP(token)
	c.pubSub.Set(otp, req)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resources.BroadcastResponse{
		OTP:            otp,
		GeneratingTime: originTime,
		Interval:       interval,
	}); err != nil {
		log.Println("error marshling send response", err)
	}
}
