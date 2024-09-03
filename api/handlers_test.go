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

		var bodyResp map[string]interface{}
		err = json.Unmarshal(response.Body.Bytes(), &bodyResp)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, "username already exists", bodyResp["error"])
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

		var bodyResp map[string]interface{}
		err = json.Unmarshal(response.Body.Bytes(), &bodyResp)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, "User added successfully", bodyResp["message"])
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

		var bodyResp map[string]interface{}
		err = json.Unmarshal(response.Body.Bytes(), &bodyResp)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, "Invalid or non-existent user", bodyResp["message"])
	})
	t.Run("Assert equal if user logged successfully", func(t *testing.T) {
		body := &models.User{
			Name:     "StoredUser",
			Password: "54321",
		}
		req, err := buildPostRequestT(body, "/signin")
		if err != nil {
			t.Fatal(err)
		}
		response := executeRequestT(req, app.router)
		checkResponseCodeT(t, http.StatusOK, response.Code)

		var bodyResp map[string]interface{}
		err = json.Unmarshal(response.Body.Bytes(), &bodyResp)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, "User logged successfully", bodyResp["message"])
	})
}
