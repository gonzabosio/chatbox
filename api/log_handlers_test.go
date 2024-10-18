package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/gonzabosio/chat-box/models"
	ws "github.com/gonzabosio/chat-box/websocket"
	"github.com/stretchr/testify/require"
)

func TestSignUp(t *testing.T) {
	app := &App{}
	if err := connectLocalMongoT(app); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(app, os.Getenv("JWT_KEY"))
	wshandler := ws.NewWSHandler(handler.service)
	app.routing(handler, wshandler)
	stored := &models.User{
		Name:     "StoredUser",
		Password: "54321",
	}
	app.client.Database("chat_box").Collection("users").InsertOne(context.TODO(), stored)
	t.Run("Assert equal if user already exists", func(t *testing.T) {
		body := &models.User{
			Name:     "StoredUser",
			Password: "54321",
		}
		req, err := buildPostRequestT(body, "/signup")
		if err != nil {
			t.Fatal(err)
		}
		response := executeRequestT(req, app.router)
		checkResponseCodeT(t, http.StatusBadRequest, response.Code)

		var respBody map[string]interface{}
		err = json.Unmarshal(response.Body.Bytes(), &respBody)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, "username already exists", respBody["error"])
	})
	t.Run("Assert equal if user is added successfully", func(t *testing.T) {
		body := &models.User{
			Name:     "NewUser",
			Password: "54321",
		}
		req, err := buildPostRequestT(body, "/signup")
		if err != nil {
			t.Fatal(err)
		}
		response := executeRequestT(req, app.router)
		checkResponseCodeT(t, http.StatusOK, response.Code)

		var respBody map[string]interface{}
		err = json.Unmarshal(response.Body.Bytes(), &respBody)
		if err != nil {
			println(respBody["error"])
			t.Fatal(err)
		}
		require.Equal(t, "User added successfully", respBody["message"])
	})
	defer app.client.Disconnect(context.TODO())
}

func TestSignIn(t *testing.T) {
	app := &App{}
	if err := connectLocalMongoT(app); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(app, os.Getenv("JWT_KEY"))
	wshandler := ws.NewWSHandler(handler.service)
	app.routing(handler, wshandler)
	t.Run("Assert equal if invalid user", func(t *testing.T) {
		body := &models.User{
			Name:     "Inexistent",
			Password: "54321",
		}
		req, err := buildPostRequestT(body, "/signin")
		if err != nil {
			t.Fatal(err)
		}
		response := executeRequestT(req, app.router)
		checkResponseCodeT(t, http.StatusUnauthorized, response.Code)

		var respBody map[string]interface{}
		err = json.Unmarshal(response.Body.Bytes(), &respBody)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, "Invalid or non-existent user", respBody["message"])
	})
	t.Run("Assert equal if user logged successfully", func(t *testing.T) {
		body := &models.User{
			Name:     "NewUser",
			Password: "54321",
		}
		req, err := buildPostRequestT(body, "/signin")
		if err != nil {
			t.Fatal(err)
		}
		response := executeRequestT(req, app.router)
		checkResponseCodeT(t, http.StatusOK, response.Code)

		var respBody map[string]interface{}
		err = json.Unmarshal(response.Body.Bytes(), &respBody)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, "User logged successfully", respBody["message"])
	})
	app.client.Database("chat_box").Collection("users").Drop(context.TODO())
}
