package repository

import (
	"github.com/gonzabosio/chat-box/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockMongoDBService struct {
	mock.Mock
}

func (m *MockMongoDBService) LoadChats(userId string) ([]models.Chat, error) {
	args := m.Called(userId)
	return args.Get(0).([]models.Chat), args.Error(1)
}
func (m *MockMongoDBService) AddChat(contact *models.Contact) (*models.Chat, error) {
	args := m.Called(contact)
	return args.Get(0).(*models.Chat), args.Error(1)
}
func (m *MockMongoDBService) DeleteChat(chatId string) error {
	args := m.Called(chatId)
	return args.Error(0)
}
func (m *MockMongoDBService) SendMessage(msgReq *models.Message) (*models.Message, error) {
	args := m.Called(msgReq)
	return args.Get(0).(*models.Message), args.Error(1)
}
func (m *MockMongoDBService) LoadMessages(chatId string) ([]models.Message, error) {
	args := m.Called(chatId)
	return args.Get(0).([]models.Message), args.Error(1)
}
func (m *MockMongoDBService) EditMessage(msgId, newMsg string) (*models.Message, error) {
	args := m.Called(msgId, newMsg)
	return args.Get(0).(*models.Message), args.Error(1)
}
func (m *MockMongoDBService) DeleteMessage(msgId string) error {
	args := m.Called(msgId)
	return args.Error(0)
}

func (m *MockMongoDBService) RegisterUser(user *models.User) (*mongo.InsertOneResult, error) {
	args := m.Called(user)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}
func (m *MockMongoDBService) LoginUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockMongoDBService) GetUserById(user *models.UserDataResponse, id primitive.ObjectID) error {
	args := m.Called(user, id)
	return args.Error(1)
}
func (m *MockMongoDBService) SaveUserPersonalDataDB(userId primitive.ObjectID, personal *models.Personal) error {
	args := m.Called(userId, personal)
	return args.Error(1)
}
func (m *MockMongoDBService) CreateSessions(session *models.Session) (*models.Session, error) {
	args := m.Called(session)
	return args.Get(0).(*models.Session), args.Error(1)
}
func (m *MockMongoDBService) GetRefresh(sessionId string) (string, error) {
	args := m.Called(sessionId)
	return args.Get(0).(string), args.Error(1)
}
func (m *MockMongoDBService) GetSessions(id string) (*models.Session, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Session), args.Error(1)
}
func (m *MockMongoDBService) RevokeSession(id string) error {
	args := m.Called(id)
	return args.Error(1)
}
func (m *MockMongoDBService) DeleteSession(id string) error {
	args := m.Called(id)
	return args.Error(1)
}
