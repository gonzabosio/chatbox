package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gonzabosio/chat-box/models"
	"github.com/gonzabosio/chat-box/repository"
	"github.com/gonzabosio/chat-box/token"
	"github.com/stretchr/testify/assert"
)

func TestLoadChat(t *testing.T) {
	// mock service and data
	chats := []models.Chat{
		{ID: "1",
			Participants: []map[string]string{
				{"id": "1", "name": "User1"},
				{"id": "2", "name": "User2"},
			}},
	}
	userID := "1234"
	mockService := new(repository.MockMongoDBService)
	mockService.On("LoadChats", userID).Return(chats, nil)

	handler := &handler{service: mockService, tokenMaker: token.NewJWTMaker("abc")}

	// build request
	req := httptest.NewRequest(http.MethodGet, "/chat/1234", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user-id", "1234")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.loadChats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, w.Body.String(), `{"chats":[{"id":"1","participants":[{"id":"1","name":"User1"},{"id":"2","name":"User2"}]}],"message":"User chats retrieved"}`)

	mockService.AssertCalled(t, "LoadChats", userID)
}

func TestAddChat(t *testing.T) {
	contact := &models.Contact{
		Username:     "User2",
		PetitionerID: "1234",
		Petitioner:   "User1",
	}
	jsonContact, err := json.Marshal(contact)
	if err != nil {
		t.Fatal("Failed JSON serialization: ", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(jsonContact))
	w := httptest.NewRecorder()

	mockService := new(repository.MockMongoDBService)
	mockService.On("AddChat", contact).Return(&models.Chat{}, nil)
	handler := &handler{service: mockService, tokenMaker: token.NewJWTMaker("abc")}

	handler.addChat(w, req)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, "Chat added successfully", resBody["message"])
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertCalled(t, "AddChat", contact)
}
func TestDeleteChat(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/chat/1234", nil)
	chatID := "1234"
	w := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("chat-id", chatID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	chi.URLParam(req, "chat-id")

	mockService := new(repository.MockMongoDBService)
	mockService.On("DeleteChat", chatID).Return(nil)
	handler := handler{service: mockService, tokenMaker: token.NewJWTMaker("abc")}
	handler.deleteChat(w, req)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, "Chat deleted successfully", resBody["message"])
	mockService.AssertCalled(t, "DeleteChat", chatID)
}
func TestLoadMessages(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/chat/{chat-id}/messages", nil)
	w := httptest.NewRecorder()
	chatID := "1234"
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("chat-id", chatID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockService := new(repository.MockMongoDBService)
	mockService.On("LoadMessages", chatID).Return([]models.Message{}, nil)

	handler := &handler{service: mockService, tokenMaker: token.NewJWTMaker("abc")}
	handler.loadMessages(w, req)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, "Messages loaded", resBody["message"])
	mockService.AssertCalled(t, "LoadMessages", chatID)
}
