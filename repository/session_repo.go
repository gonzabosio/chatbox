package repository

import (
	"context"
	"fmt"

	"github.com/gonzabosio/chat-box/models"
	"go.mongodb.org/mongo-driver/bson"
)

type SessionsRepository interface {
	CreateSessions(session *models.Session) (*models.Session, error)
	GetRefresh(sessionId string) (string, error)
	GetSessions(id string) (*models.Session, error)
	RevokeSession(id string) error
	DeleteSession(id string) error
}

func (ms *MongoDBService) CreateSessions(session *models.Session) (*models.Session, error) {
	coll := ms.DB.Collection("sessions")
	result, err := coll.InsertOne(context.TODO(), session)
	if err != nil {
		return nil, fmt.Errorf("error creating session: %v", err)
	}
	fmt.Println(result.InsertedID)
	return session, nil
}

func (ms *MongoDBService) GetRefresh(sessionId string) (string, error) {
	var res models.Session
	filter := bson.D{{Key: "_id", Value: sessionId}}
	err := ms.DB.Collection("sessions").FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return "", err
	}
	return res.RefreshToken, nil
}

func (ms *MongoDBService) GetSessions(id string) (*models.Session, error) {
	var s *models.Session
	filter := bson.D{{Key: "_id", Value: id}}
	err := ms.DB.Collection("sessions").FindOne(context.TODO(), filter).Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %v", err)
	}
	return s, nil
}

func (ms *MongoDBService) RevokeSession(id string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_revoked", Value: true}}}}
	err := ms.DB.Collection("sessions").FindOneAndUpdate(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error revoking session: %v", err)
	}
	return nil
}

func (ms *MongoDBService) DeleteSession(id string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	err := ms.DB.Collection("sessions").FindOneAndDelete(context.TODO(), filter).Err()
	if err != nil {
		return fmt.Errorf("error deleting session: %v", err)
	}
	return nil
}
