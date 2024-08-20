package repo

import (
	"context"
	"fmt"

	"github.com/gonzabosio/chat-box/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type MongoDBService struct {
	DB *mongo.Database
}

type UserRepository interface {
	RegisterUser(user *storage.User) (*mongo.InsertOneResult, error)
	LoginUser(user *storage.User) (*storage.User, error)
	GetUserById(user *storage.User, filter primitive.D) error
}

func comparePasswords(dbPassword, password string) bool {
	fmt.Println(dbPassword)
	fmt.Println(password)
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	return err == nil
}

func (ms *MongoDBService) RegisterUser(user *storage.User) (*mongo.InsertOneResult, error) {

	filter := bson.D{{Key: "username", Value: user.Username}}
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

func (ms *MongoDBService) LoginUser(user *storage.User) (*storage.User, error) {
	filter := bson.D{{Key: "username", Value: user.Username}}
	coll := ms.DB.Collection("users")
	var dbUser *storage.User
	err := coll.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {
		return nil, err
	}
	if !comparePasswords(dbUser.Password, user.Password) {
		return nil, fmt.Errorf("incorrect username or password")
	}
	return dbUser, nil
}

func (ms *MongoDBService) GetUserById(user *storage.User, filter primitive.D) error {
	coll := ms.DB.Collection("users")
	if err := coll.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return err
	}
	return nil
}
