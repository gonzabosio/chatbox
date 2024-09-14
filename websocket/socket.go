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
		return true
	},
}

var clientsToSend = make(map[*websocket.Conn]bool)
var broadcastToSend = make(chan *models.Message)

func (h *WSHandler) SendMsgWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not upgrade to websocket connection: ", err)
		return
	}
	defer c.Close()
	// Add a new connected client
	clientsToSend[c] = true

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			delete(clientsToSend, c)
			c.WriteJSON(map[string]string{
				"message": "Could not read message to send",
				"error":   err.Error(),
			})
			break
		}

		var body *models.Message
		if err = json.Unmarshal(msg, &body); err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not unmarshal request body",
				"error":   err.Error(),
			})
			break
		}

		newMsg, err := h.service.SendMessages(body)
		if err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not save message in database",
				"error":   err.Error(),
			})
			break
		}
		// Call response handler
		broadcastToSend <- newMsg
	}
}

type editMsgReq struct {
	MessageID  string `json:"message_id"`
	NewMessage string `json:"new_message"`
}

var clientsToEdit = make(map[*websocket.Conn]bool)
var broadcastToEdit = make(chan *models.Message)

func (h *WSHandler) EditMsgWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not upgrade to ws prot: ", err)
		return
	}
	defer c.Close()
	clientsToEdit[c] = true

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			delete(clientsToEdit, c)
			c.WriteJSON(map[string]string{
				"message": "Could not read message update request",
				"error":   err.Error(),
			})
			break
		}

		var body *editMsgReq
		if err = json.Unmarshal(msg, &body); err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not unmarshal request body",
				"error":   err.Error(),
			})
			break
		}

		newMsg, err := h.service.EditMessage(body.MessageID, body.NewMessage)
		if err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not edit message",
				"error":   err.Error(),
			})
			break
		}
		broadcastToEdit <- newMsg
	}
}

var clientsToDelete = make(map[*websocket.Conn]bool)
var broadcastToDelete = make(chan string)

func (h *WSHandler) DeleteMsgWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not upgrade to ws prot: ", err)
		return
	}
	defer c.Close()
	clientsToDelete[c] = true

	for {
		_, msgID, err := c.ReadMessage()
		if err != nil {
			delete(clientsToDelete, c)
			c.WriteJSON(map[string]string{
				"message": "Bad request to delete message",
				"error":   err.Error(),
			})
			break
		}
		id := string(msgID)
		err = h.service.DeleteMessage(id)
		if err != nil {
			c.WriteJSON(map[string]string{
				"message": "Could not delete message from db",
				"error":   err.Error(),
			})
			break
		}
		broadcastToDelete <- id
	}
}

func HandleWebSocketSender() {
	for {
		newMsg := <-broadcastToSend
		for client := range clientsToSend {
			if err := client.WriteJSON(&newMsg); err != nil {
				log.Println("Error writing websocket response: ", err)
				client.WriteJSON(map[string]string{
					"message": err.Error(),
				})
				delete(clientsToSend, client)
				client.Close()
			}
		}
	}
}

func HandleWebSocketEditor() {
	for {
		newMsg := <-broadcastToEdit
		for client := range clientsToEdit {
			if err := client.WriteJSON(&newMsg); err != nil {
				log.Println("Error writing websocket response: ", err)
				client.WriteJSON(map[string]string{
					"message": err.Error(),
				})
				delete(clientsToEdit, client)
				client.Close()
			}
		}
	}
}

func HandleWebSocketDelete() {
	for {
		id := <-broadcastToDelete
		for client := range clientsToDelete {
			if err := client.WriteJSON(map[string]string{
				"message":    "Message deleted successfully",
				"message_id": id,
			}); err != nil {
				log.Println("Error writing websocket response:", err)
				client.WriteJSON(map[string]string{
					"message": err.Error(),
				})
				delete(clientsToDelete, client)
				client.Close()
			}
		}
	}
}
