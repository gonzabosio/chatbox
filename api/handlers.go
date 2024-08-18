package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gonzabosio/chat-box/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func (a *App) dbInstance() {
	db = a.db
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func register(w http.ResponseWriter, r *http.Request) {
	var user storage.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("Bad request to register user: %v", err),
		})
		return
	}
	coll := db.Collection("users")
	res, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf("Could not post the user in the database: %v", err),
		})
		return
	}

	token, err := generateJWT(res.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf("Could not generate the JWT: %v", err),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User added successfully",
		"user_id": res.InsertedID,
		"token":   token,
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	filter := bson.D{{Key: "username", Value: username}, {Key: "password", Value: password}}
	dbUser := new(storage.User)
	coll := db.Collection("users")
	err := coll.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{
			"message": fmt.Sprintf("Invalid or non-existent user | %v", err),
		})
		return
	}
	token, err := generateJWT(dbUser.Id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf("Could not generate the JWT: %v", err),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User logged successfully",
		"user_id": dbUser.Id,
		"token":   token,
	})
}

func getUserDataById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	dbUser := new(storage.User)
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "could not convert id to ObjectId",
		})
	}
	filter := bson.D{{Key: "_id", Value: id}}
	coll := db.Collection("users")
	if err := coll.FindOne(context.TODO(), filter).Decode(&dbUser); err != nil {
		respondJSON(w, http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Invalid or non-existent %v", err),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":   "User data retrieved",
		"user_data": dbUser,
		"claims": map[string]any{
			"subject":    r.Context().Value(ctxKey(userId)).(string),
			"expires_at": r.Context().Value(ctxKey(expiresAt)).(string),
		},
	})
}
