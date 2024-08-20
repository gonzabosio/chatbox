package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gonzabosio/chat-box/repo"
	"github.com/gonzabosio/chat-box/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// function to instance service is in common.go
var ms *repo.MongoDBService

func signUp(w http.ResponseWriter, r *http.Request) {
	var user storage.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Bad request to register user",
			"error":   err.Error(),
		})
		return
	}
	hashedPassw, err := hashPassword(user.Password)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Could not hash password",
			"error":   err.Error(),
		})
		return
	}
	user.Password = hashedPassw
	res, err := ms.RegisterUser(&user)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Could not register the user in the database",
			"error":   err.Error(),
		})
		return
	}
	token, err := generateJWT(res.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not generate the JWT",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User added successfully",
		"user_id": res.InsertedID,
		"token":   token,
	})
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var user storage.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Bad request to login user",
			"error":   err.Error(),
		})
		return
	}
	dbUser, err := ms.LoginUser(&user)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{
			"message": "Invalid or non-existent user",
			"error":   err.Error(),
		})
		return
	}
	token, err := generateJWT(dbUser.Id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not generate the JWT",
			"error":   err.Error(),
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
			"message": "Could not convert id to ObjectId",
		})
		return
	}
	filter := bson.D{{Key: "_id", Value: id}}
	if err := ms.GetUserById(dbUser, filter); err != nil {
		respondJSON(w, http.StatusOK, map[string]string{
			"message": "Invalid or non-existent user id",
			"error":   err.Error(),
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
