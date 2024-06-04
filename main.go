package main

import (
	"goshare/internal/config"
	"goshare/internal/otp"
	"goshare/internal/server"
	"goshare/internal/store"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	otpGen := otp.NewTotp()
	store := store.NewStore(&otpGen)
	s := server.NewServer(cfg, store)

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
