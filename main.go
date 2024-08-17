package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gonzabosio/chat-box/api"
)

func main() {
	app := &api.App{}
	if err := app.InitServer(); err != nil {
		log.Fatalf("Could not start server instance: %v", err)
	}

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("Server run failed: %v", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh)
	<-signalCh
	app.ShutdownConn()
}
