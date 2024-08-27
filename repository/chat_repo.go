package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gonzabosio/chat-box/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository interface {
	LoadChats(userId string) ([]models.Chat, error)
	AddContact(contact *models.Contact) error
	DeleteContact(chatId string) error
	SendMessages(msgReq *models.Message) error
	LoadMessages(chatId string) ([]models.Message, error)
	EditMessage(msgId string) error
	DeleteMessage(msgId string) error
}

// mongodb instance in user_repo.go
func (ms *MongoDBService) LoadChats(userId string) (chats []models.Chat, err error) {
	coll := ms.DB.Collection("chats")
	filter := bson.M{
		"participants": bson.M{
			"$elemMatch": bson.M{
				"id": userId,
			},
		},
	}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &chats)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (ms *MongoDBService) AddChat(contact *models.Contact) (interface{}, error) {
	var newContact models.User
	coll := ms.DB.Collection("users")
	filter := bson.D{{Key: "name", Value: contact.Username}}
	err := coll.FindOne(context.TODO(), filter).Decode(&newContact)
	if err != nil {
		return nil, fmt.Errorf("could not found the contact: %v", err)
	}

	coll = ms.DB.Collection("chats")
	filterAlreadyAdded := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "participants", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "id", Value: newContact.ID}}}}}},
			bson.D{{Key: "participants", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "id", Value: contact.PetitionerID}}}}}},
		}},
	}
	result := coll.FindOne(context.TODO(), filterAlreadyAdded)
	if result.Err() != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("contact is already in the chats")
	}

	var chat models.Chat
	chat.Participants = append(chat.Participants, map[string]string{
		"id":   newContact.ID,
		"name": contact.Username,
	})
	chat.Participants = append(chat.Participants, map[string]string{
		"id":   contact.PetitionerID,
		"name": contact.Petitioner,
	})
	res, err := coll.InsertOne(context.TODO(), chat)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (ms *MongoDBService) DeleteChat(chatId string) error {
	coll := ms.DB.Collection("chats")
	id, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MongoDBService) LoadMessages(chatId string) ([]models.Message, error) {
	coll := ms.DB.Collection("messages")
	var mesagges []models.Message
	filter := bson.D{{Key: "chat_id", Value: chatId}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &mesagges)
	if err != nil {
		return nil, err
	}
	return mesagges, nil
}

func (ms *MongoDBService) SendMessages(msgReq *models.Message) error {
	msgReq.SentAt = time.Now()
	coll := ms.DB.Collection("messages")
	_, err := coll.InsertOne(context.TODO(), msgReq)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MongoDBService) EditMessage(msgId, newMsg string) error {
	coll := ms.DB.Collection("messages")
	id, err := primitive.ObjectIDFromHex(msgId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "content", Value: newMsg}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MongoDBService) DeleteMessage(msgId string) error {
	coll := ms.DB.Collection("messages")
	id, err := primitive.ObjectIDFromHex(msgId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
