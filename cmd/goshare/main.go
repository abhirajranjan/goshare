package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goshare/internal/handler"
	"goshare/internal/pubsub"
	"goshare/internal/server"
)

func main() {
	adapter := pubsub.NewPubsub()
	srvhandler := handler.NewHandler(adapter)
	httpsrv := &http.Server{
		Addr:    ":8080",
		Handler: srvhandler.Router,
	}
	srv := server.NewServer(httpsrv)

	log.Println("starting server at ", httpsrv.Addr)
	srv.Start()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	t := time.NewTimer(time.Second)

Termination:
	for {
		select {
		// so this go routine doesnt block,
		// if all go routine will be blocked
		// then program terminates from deadlock
		case <-t.C:
			t.Reset(time.Second)
		case s := <-sigc:
			log.Printf("signal %s encountered. Terminating...\n", s)
			srv.Close()
		case <-srv.Wait:
			break Termination
		}
	}

	if srv.Error() != nil {
		fmt.Println(srv.Error())
	}
}
