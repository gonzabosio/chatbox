package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gonzabosio/chat-box/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// response json structure
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// testing resources (func functionNameT() {})
func executeRequestT(req *http.Request, router *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCodeT(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d. got %d\n", expected, actual)
	}
}

func connectLocalMongoT(app *App) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}
	app.client = client
	fmt.Println("MongoDB Connected")
	return nil
}

func buildPostRequestT(body *models.User, endpoint string) (*http.Request, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	return req, nil
}
