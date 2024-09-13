package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gonzabosio/chat-box/api"
	ws "github.com/gonzabosio/chat-box/websocket"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Loading enviroment variables failed: %v", err)
	}

	app := &api.App{}
	if err := app.InitServer(); err != nil {
		log.Fatalf("Could not start server instance: %v", err)
	}

	go ws.HandleWebSocketSender()
	go ws.HandleWebSocketEditor()

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("Server run failed: %v", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	app.ShutdownConn()
}
