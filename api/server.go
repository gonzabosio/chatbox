package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	router *chi.Mux
	client *mongo.Client
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
		return err
	}
	// Confirm successful connection sending a ping
	if err := a.client.Database("chat_box").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}
	fmt.Println("Successfully connected to MongoDB!")

	a.router = chi.NewRouter()
	a.router.Post("/login", login)
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
