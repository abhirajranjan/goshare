package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

type otpGen interface {
	GenerateOTP(secret string) (token string, origin time.Time, interval time.Duration)
}

type pubSub interface {
	Set(id string, data any)
	Get(id string) (any, error)
}

type crypto interface {
	Sign(deviceId []byte) ([]byte, error)
	Verify(deviceId []byte) bool
	GetMessage(token []byte) []byte
}

type Controller struct {
	pubSub pubSub
	otpGen otpGen
	crypto crypto
}

func NewController(pubsub pubSub, otpgen otpGen, crypto crypto) *Controller {
	return &Controller{
		pubSub: pubsub,
		otpGen: otpgen,
		crypto: crypto,
	}
}

func extractDeviceToken(r *http.Request, c crypto) (string, error) {
	auth := r.Header.Get("Authorization")
	_, bearerToken, found := strings.Cut(auth, "Bearer ")
	if !found {
		return "", errors.New("missing deviceId")
	}

	if ok := c.Verify([]byte(bearerToken)); !ok {
		return "", errors.New("invalid deviceID")
	}

	return string(c.GetMessage([]byte(bearerToken))), nil
}
