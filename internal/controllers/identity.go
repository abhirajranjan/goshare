package controllers

import (
	"fmt"
	r "math/rand/v2"
	"net/http"
)

func (c Controller) IdentityHandler(w http.ResponseWriter, r *http.Request) {
	deviceName := randomName(6)
	fmt.Println(deviceName)
	deviceNameSigned, err := c.crypto.Sign(deviceName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error signing device"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(deviceNameSigned)
}

func randomName(length int) []byte {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	buf := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		buf[i] = alphabet[r.Int32N(26)]
	}

	return buf
}
