package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gonzabosio/chat-box/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// HANDLERS RESOURCES
// response json structure
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// JWT generation used in signUp and signIn handlers
func generateJWT(id string) (string, error) {
	// Set custom claims and create token
	claims := jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Encode token
	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return t, nil
}

// TESTING RESOURCES
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

func buildPostRequestT(body *storage.User, endpoint string) (*http.Request, error) {
	bodyReq, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(bodyReq))
	if err != nil {
		return nil, err
	}
	return req, nil
}
