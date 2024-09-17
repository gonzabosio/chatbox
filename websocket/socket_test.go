package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gonzabosio/chat-box/models"
	"github.com/gonzabosio/chat-box/repository"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	go HandleWebSocketSender()
	message := &models.Message{
		ChatID:   "1",
		SenderID: "1234",
		Content:  "Hello",
	}
	reqBody, err := json.Marshal(message)
	if err != nil {
		t.Fatal("could not marshall message: ", err)
	}

	// set up mock service
	mockService := new(repository.MockMongoDBService)
	mockService.On("SendMessage", message).Return(message, nil)

	// set up websocket handler and server
	wsHandler := NewWSHandler(mockService)
	server := httptest.NewServer(http.HandlerFunc(wsHandler.SendMsgWS))
	u := "ws" + server.URL[len("http"):]
	wsURL, err := url.Parse(u)
	if err != nil {
		t.Fatalf("could not parse url: %d", err)
	}
	defer server.Close()
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Fatalf("cannot make websocket connection: %v", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, reqBody)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}
	_, res, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}
	mockService.AssertCalled(t, "SendMessage", message)

	var response map[string]interface{}
	err = json.Unmarshal(res, &response)
	if err != nil {
		t.Fatalf("could not unmarshall response: %v", err)
	}
	assert.Equal(t, "Message was sent successfully", response["message"])
}

func TestEditMessage(t *testing.T) {
	go HandleWebSocketEditor()
	msgID := "54321"
	request := &editMsgReq{
		MessageID:  msgID,
		NewMessage: "Bye",
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		t.Fatal("could not marshall message: ", err)
	}

	mockService := new(repository.MockMongoDBService)
	mockService.On("EditMessage", msgID, request.NewMessage).Return(&models.Message{}, nil)

	wsHandler := NewWSHandler(mockService)
	server := httptest.NewServer(http.HandlerFunc(wsHandler.EditMsgWS))
	u := "ws" + server.URL[len("http"):]
	wsURL, err := url.Parse(u)
	if err != nil {
		t.Fatalf("could not parse url: %d", err)
	}
	defer server.Close()
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Fatalf("cannot make websocket connection: %v", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, reqBody)
	if err != nil {
		t.Fatalf("failed to send message to edit: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(msg, &response)
	if err != nil {
		t.Fatal("could not unmarshall response")
	}
	mockService.AssertCalled(t, "EditMessage", msgID, request.NewMessage)
	assert.Equal(t, "Message edited successfully", response["message"])
}
func TestDeleteMessage(t *testing.T) {
	go HandleWebSocketDelete()
	msgID := "1234"
	mockService := new(repository.MockMongoDBService)
	mockService.On("DeleteMessage", msgID).Return(nil)

	wsHandler := NewWSHandler(mockService)
	server := httptest.NewServer(http.HandlerFunc(wsHandler.DeleteMsgWS))
	u := "ws" + server.URL[len("http"):]
	wsURL, err := url.Parse(u)
	if err != nil {
		t.Fatalf("could not parse url: %d", err)
	}
	defer server.Close()
	time.Sleep(500 * time.Millisecond)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Fatalf("cannot make websocket connection: %v", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("1234"))
	if err != nil {
		t.Fatalf("failed to send message to edit: %v", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}
	mockService.AssertCalled(t, "DeleteMessage", msgID)
	var response map[string]interface{}
	err = json.Unmarshal(msg, &response)
	if err != nil {
		t.Fatal("could not unmarshall response")
	}

	assert.Equal(t, "Message deleted successfully", response["message"])
}
