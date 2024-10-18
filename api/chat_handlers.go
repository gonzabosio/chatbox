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
		respondJSON(w, http.StatusBadRequest, map[string]string{
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
	if err := validate.Struct(contact); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Bad request to add contact",
			"error":   err.Error(),
		})
		return
	}
	newChat, err := h.service.AddChat(&contact)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
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
		respondJSON(w, http.StatusBadRequest, map[string]string{
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
		respondJSON(w, http.StatusBadRequest, map[string]string{
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
