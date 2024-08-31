package models

import "time"

type Session struct {
	ID           string    `bson:"_id" json:"id"`
	Username     string    `bson:"username" json:"username"`
	RefreshToken string    `bson:"refresh_token" json:"refresh_token"`
	IsRevoked    bool      `bson:"is_revoked" json:"is_revoked"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	ExpiresAt    time.Time `bson:"expires_at" json:"expires_at"`
}

type SessionReq struct {
	ID string `bson:"_id" json:"session_id"`
}
