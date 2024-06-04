///$(which go) run $0 $@; exit $?

package main

import (
	"net/http"

	"goshare/config"
	"goshare/otp"
	"goshare/server"
	"goshare/store"
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
