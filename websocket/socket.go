package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gonzabosio/chat-box/models"
	"github.com/gonzabosio/chat-box/repository"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	service *repository.MongoDBService
}

func NewWSHandler(repo *repository.MongoDBService) *WSHandler {
	return &WSHandler{service: repo}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		return true
	},
}

func (h *WSHandler) ChatMsgHandler(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not upgrade to ws prot: ", err)
		return
	}
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			// show error in body response
		}
		log.Println(string(msg))

		var body *models.Message
		if err = json.Unmarshal(msg, &body); err != nil {
			log.Println(err)
			continue
		}
		log.Println(body)

		h.service.SendMessages(body)

		if err = c.WriteJSON(&body); err != nil {
			log.Println(err)
			// show error in body response
		}
	}
}
