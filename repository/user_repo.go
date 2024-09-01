package repository

import (
	"context"
	"fmt"

	"github.com/gonzabosio/chat-box/models"
	"github.com/gonzabosio/chat-box/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBService struct {
	DB *mongo.Database
}

type UserRepository interface {
	RegisterUser(user *models.User) (*mongo.InsertOneResult, error)
	LoginUser(user *models.User) (*models.User, error)
	GetUserById(user *models.User, id primitive.ObjectID) error
}

func (ms *MongoDBService) RegisterUser(user *models.User) (*mongo.InsertOneResult, error) {

	filter := bson.D{{Key: "name", Value: user.Name}}
	oneRes := ms.DB.Collection("users").FindOne(context.TODO(), filter)

	if oneRes.Decode(user) != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("username already exists")
	}
	coll := ms.DB.Collection("users")
	addedUser, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return addedUser, nil
}

func (ms *MongoDBService) LoginUser(user *models.User) (*models.User, error) {
	filter := bson.D{{Key: "name", Value: user.Name}}
	coll := ms.DB.Collection("users")
	var dbUser *models.User
	err := coll.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {
		return nil, err
	}
	if !utils.ComparePasswords(dbUser.Password, user.Password) {
		return nil, fmt.Errorf("incorrect username or password")
	}
	return dbUser, nil
}

func (ms *MongoDBService) GetUserById(user *models.User, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	coll := ms.DB.Collection("users")
	if err := coll.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return err
	}
	return nil
}
