package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	ws "github.com/gonzabosio/chat-box/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	router *chi.Mux
	client *mongo.Client
}

func (a *App) InitServer() error {
	// New client and server connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("ATLAS_URI")).SetServerAPIOptions(serverAPI)
	ctx, cancelCtx := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelCtx()
	var err error
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

	handler := NewHandler(a, os.Getenv("JWT_KEY"))
	wshandler := ws.NewWSHandler(handler.service)
	a.routing(handler, wshandler)

	return nil
}

func (a *App) routing(h *handler, wsh *ws.WSHandler) {
	a.router = chi.NewRouter()
	a.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("FRONT_URL")},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         1800,
	}))
	//Public
	a.router.Use(middleware.Logger)
	a.router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(10, 1*time.Minute))
		r.Post("/signup", h.signUp)
		r.Post("/signin", h.signIn)
		r.Delete("/logout/{session-id}", h.logout)
		a.router.Route("/token", func(r chi.Router) {
			r.Post("/renew/{session-id}", h.renewAccessToken)
			r.Post("/revoke/{session-id}", h.revokeSession)
		})
	})
	//Private
	a.router.Group(func(r chi.Router) {
		r.Use(h.authMiddleware)
		r.Use(httprate.LimitByIP(50, 1*time.Minute))
		r.Route("/user", func(r chi.Router) {
			r.Get("/{id}", h.getUserDataById)
		})
		r.Route("/chat", func(r chi.Router) {
			r.Get("/load/{user-id}", h.loadChats)
			r.Post("/add", h.addChat)
			r.Delete("/delete/{chat-id}", h.deleteChat)
			r.Get("/load-messages/{chat-id}", h.loadMessages)
			r.Delete("/delete-message/{msg-id}", h.deleteMessage)
		})
		a.router.Route("/ws", func(r chi.Router) {
			r.HandleFunc("/send-msg", wsh.SendMsgWS)
			r.HandleFunc("/edit-msg", wsh.EditMsgWS)
		})
	})
}

func (a *App) Run() error {
	return http.ListenAndServe(":"+os.Getenv("PORT"), a.router)
}

func (a *App) ShutdownConn() {
	if err := a.client.Disconnect(context.TODO()); err != nil {
		fmt.Printf("Mongo disconnection failed: %v", err)
	} else {
		fmt.Println("MongoDB Disconnected")
	}
}
