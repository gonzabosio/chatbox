package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gonzabosio/chat-box/models"
)

// handler struct in log_handlers.go
func (h *handler) loadChats(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user-id")
	chats, err := h.service.LoadChats(userId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not load user chats",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User chats retrieved",
		"chats":   chats,
	})
}

func (h *handler) addChat(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid contact data",
			"error":   err.Error(),
		})
		return
	}
	newChat, err := h.service.AddChat(&contact)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Chat added successfully",
		"chat":    newChat,
		"contact": contact,
	})
}

func (h *handler) deleteChat(w http.ResponseWriter, r *http.Request) {
	chatId := chi.URLParam(r, "chat-id")
	err := h.service.DeleteChat(chatId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not delete chat",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Chat deleted successfully"})
}

func (h *handler) loadMessages(w http.ResponseWriter, r *http.Request) {
	chatId := chi.URLParam(r, "chat-id")
	messages, err := h.service.LoadMessages(chatId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not load messages",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Messages loaded",
		"messages": messages,
	})
}

func (h *handler) sendMessage(w http.ResponseWriter, r *http.Request) {
	bodyReq := new(models.Message)
	err := json.NewDecoder(r.Body).Decode(&bodyReq)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Could not read body request to send message",
			"error":   err.Error(),
		})
		return
	}
	err = h.service.SendMessages(bodyReq)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not send the message",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Message was sent",
		"doc":     bodyReq,
	})
}

func (h *handler) editMessage(w http.ResponseWriter, r *http.Request) {
	type EditBodyReq struct {
		MessageID  string `json:"message_id"`
		NewMessage string `json:"new_message"`
	}
	bodyReq := new(EditBodyReq)
	err := json.NewDecoder(r.Body).Decode(&bodyReq)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Could not read body request to edit the message",
			"error":   err.Error(),
		})
		return
	}
	newMsg, err := h.service.EditMessage(bodyReq.MessageID, bodyReq.NewMessage)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not edit message",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Message edited",
		"new_message": newMsg,
	})
}

func (h *handler) deleteMessage(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "msg-id")
	err := h.service.DeleteMessage(msgId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not delete message",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Message deleted"})
}
