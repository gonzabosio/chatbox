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

func (h *WSHandler) SendMsgWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not upgrade to ws prot: ", err)
		return
	}
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not read message to send",
				"error":   err.Error(),
			})
		}
		log.Println(string(msg))

		var body *models.Message
		if err = json.Unmarshal(msg, &body); err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not unmarshal request body",
				"error":   err.Error(),
			})
		}
		log.Println(body)

		newMsg, err := h.service.SendMessages(body)
		if err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not respond save message in database",
				"error":   err.Error(),
			})
		}

		if err = c.WriteJSON(&newMsg); err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not respond with sent message",
				"error":   err.Error(),
			})
		}
	}
}

func (h *WSHandler) EditMsgWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not upgrade to ws prot: ", err)
		return
	}
	defer c.Close()
	type editMsgReq struct {
		MessageID  string `json:"message_id"`
		NewMessage string `json:"new_message"`
	}
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not read message update request",
				"error":   err.Error(),
			})
		}
		log.Println(string(msg))

		var body *editMsgReq
		if err = json.Unmarshal(msg, &body); err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not unmarshal request body",
				"error":   err.Error(),
			})
		}
		log.Println(body)

		newMsg, err := h.service.EditMessage(body.MessageID, body.NewMessage)
		if err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not edit message",
				"error":   err.Error(),
			})
		}
		if err = c.WriteJSON(&newMsg); err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not respond with updated message",
				"error":   err.Error(),
			})
		}
	}
}
