package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gonzabosio/chat-box/models"
	"github.com/gonzabosio/chat-box/repository"
	"github.com/gonzabosio/chat-box/token"
	"github.com/gonzabosio/chat-box/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	service    *repository.MongoDBService
	tokenMaker *token.JWTMaker
}

func NewHandler(app *App, secretKey string) *handler {
	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return &handler{
		service:    &repository.MongoDBService{DB: app.client.Database("chat_box")},
		tokenMaker: token.NewJWTMaker(string(key)),
	}
}

func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Bad request to register user",
			"error":   err.Error(),
		})
		return
	}
	if err := checkLoginValues(user.Name, user.Password); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Bad request to register user",
			"error":   err.Error(),
		})
		return
	}
	hashedPassw, err := utils.HashPassword(user.Password)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Could not hash password",
			"error":   err.Error(),
		})
		return
	}
	user.Password = hashedPassw
	res, err := h.service.RegisterUser(&user)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Could not register the user in the database",
			"error":   err.Error(),
		})
		return
	}
	accessToken, accessClaims, err := h.tokenMaker.CreateToken(user.ID, user.Name, 15*time.Minute)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not generate the access JWT",
			"error":   err.Error(),
		})
		return
	}
	refreshToken, refreshClaims, err := h.tokenMaker.CreateToken(user.ID, user.Name, 24*time.Hour)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not generate the refresh JWT",
			"error":   err.Error(),
		})
		return
	}
	session, err := h.service.CreateSessions(&models.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		Username:     user.Name,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not create session",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "User added successfully",
		"session_id":   session.ID,
		"access_token": accessToken,
		"access_exp":   accessClaims.ExpiresAt.Time,
		"refresh_exp":  refreshClaims.ExpiresAt.Time,
		"user": map[string]interface{}{
			"id":   res.InsertedID,
			"name": user.Name,
		},
	})
}

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Bad request to login user",
			"error":   err.Error(),
		})
		return
	}
	if err := checkLoginValues(user.Name, user.Password); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Bad request to login user",
			"error":   err.Error(),
		})
		return
	}
	dbUser, err := h.service.LoginUser(&user)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{
			"message": "Invalid or non-existent user",
			"error":   err.Error(),
		})
		return
	}
	accessToken, accessClaims, err := h.tokenMaker.CreateToken(dbUser.ID, dbUser.Name, 15*time.Minute)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not generate the access JWT",
			"error":   err.Error(),
		})
		return
	}
	refreshToken, refreshClaims, err := h.tokenMaker.CreateToken(dbUser.ID, dbUser.Name, 24*time.Hour)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not generate the refresh JWT",
			"error":   err.Error(),
		})
		return
	}
	session, err := h.service.CreateSessions(&models.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		Username:     dbUser.Name,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not create session",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "User logged successfully",
		"session_id":   session.ID,
		"access_token": accessToken,
		"access_exp":   accessClaims.ExpiresAt.Time,
		"refresh_exp":  refreshClaims.ExpiresAt.Time,
		"user": map[string]interface{}{
			"id":   dbUser.ID,
			"name": dbUser.Name,
		},
	})
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "session-id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Missing session id"})
		return
	}
	err := h.service.DeleteSession(id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Could not delete session",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "User logged out"})
}

func (h *handler) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	sessionId := chi.URLParam(r, "session-id")
	refreshTokenStr, err := h.service.GetRefresh(sessionId)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Could not found user session",
			"error":   err.Error(),
		})
		return
	}
	refreshClaims, err := h.tokenMaker.VerifyToken(refreshTokenStr)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{
			"message": "Error verifying token",
			"error":   err.Error(),
		})
		return
	}
	session, err := h.service.GetSessions(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	if session.IsRevoked {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "session revoked"})
		return
	}
	if session.Username != refreshClaims.Username {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid session"})
		return
	}
	accessToken, accessClaims, err := h.tokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Username, 15*time.Minute)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "error creating token",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "Access token renewed",
		"access_token": accessToken,
		"access_exp":   accessClaims.ExpiresAt.Time,
	})
}

func (h *handler) revokeSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "session-id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Missing session id"})
		return
	}
	err := h.service.RevokeSession(id)
	if err != nil {
		if id == "" {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error revoking session"})
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) getUserDataById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	dbUser := new(models.User)
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Could not convert id to ObjectId"})
		return
	}
	if err := h.service.GetUserById(dbUser, id); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Invalid or non-existent user id",
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":   "User data retrieved",
		"user_data": dbUser,
		"claims": map[string]any{
			"subject":    r.Context().Value(ctxKey(claimsId)).(string),
			"expires_at": r.Context().Value(ctxKey(expiresAt)).(string),
		},
	})
}
