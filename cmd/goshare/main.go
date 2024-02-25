package main

import (
	"goshare/internal/config"
	"goshare/internal/controllers"
	"goshare/internal/crypto"
	"goshare/internal/server"
	"goshare/pkg/otp"
	"goshare/pkg/pubsub"
	"log"
)

func main() {
	cfg := config.GetConfig()

	ps := pubsub.NewPubSub(cfg.Pubsub.TTL)
	otp := otp.NewTOTP(otp.WithDuration(cfg.Otp.ValidDuration))
	crypt := crypto.NewCrypto(cfg.Key.PublicKey, cfg.Key.PrivateKey)

	ctrl := controllers.NewController(ps, &otp, &crypt)
	srv := server.NewServer(ctrl)

	log.Printf("starting server at %s", cfg.Server.Addr)
	log.Printf("server: %s\n", srv.Start(cfg.Server.Addr))
}
