package repository

import "go.mongodb.org/mongo-driver/mongo"

type Service interface {
	UserRepository
	ChatRepository
	SessionsRepository
}

type MongoDBService struct {
	DB *mongo.Database
}
