package models

import "time"

type Chat struct {
	ID           string            `bson:"_id,omitempty" json:"id,omitempty"`
	Participants map[string]string `bson:"participants" json:"participants"`
}

type Message struct {
	ID       string    `bson:"_id,omitempty" json:"id,omitempty"`
	ChatID   string    `bson:"chat_id" json:"chat_id"`
	SenderID string    `bson:"sender_id" json:"sender_id"`
	Content  string    `bson:"content" json:"content"`
	SentAt   time.Time `bson:"sent_at" json:"sent_at,omitempty"`
}

type Contact struct {
	ID           string `json:"contact_id,omitempty"`
	Username     string `json:"username"`
	PetitionerID string `json:"petitioner_id"`
	Petitioner   string `json:"petitioner"`
}
