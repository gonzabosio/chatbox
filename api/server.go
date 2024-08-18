package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	router *chi.Mux
	client *mongo.Client
	db     *mongo.Database
}

func (a *App) InitServer() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	// New client and server connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("ATLAS_URI")).SetServerAPIOptions(serverAPI)
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()
	a.client, err = mongo.Connect(ctx, opts)
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("timed out connecting MongoDB")
	} else if ctx.Err() == context.Canceled {
		return fmt.Errorf("connection to MongoDB canceled")
	}
	if err != nil {
		return fmt.Errorf("connection to MongoDB failed | %v", err)
	}
	// Confirm successful connection sending a ping
	if err := a.client.Database("chat_box").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}
	fmt.Println("Successfully connected to MongoDB!")
	a.db = a.client.Database("chat_box")
	a.dbInstance()

	a.router = chi.NewRouter()

	a.router.Use(middleware.Logger)
	a.router.Group(func(r chi.Router) {
		r.Post("/register", register)
		r.Post("/login", login)
	})
	a.router.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/user/{id}", getUserDataById)
	})
	return nil
}

func (a *App) Run() error {
	return http.ListenAndServe(":8000", a.router)
}

func (a *App) ShutdownConn() {
	if err := a.client.Disconnect(context.TODO()); err != nil {
		fmt.Printf("Mongo disconnection failed: %v", err)
	} else {
		fmt.Println("MongoDB Disconnected")
	}
}

func generateJWT(id string) (string, error) {
	// Set custom claims and create token
	claims := jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 60)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Encode token
	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return t, nil
}
