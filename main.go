package main

import (
	"log"

	"github.com/gonzabosio/chat-box/api"
)

func main() {
	app := &api.App{}
	app.InitServer()

	if err := app.Run(); err != nil {
		log.Fatalf("Server run failed")
	}
}
