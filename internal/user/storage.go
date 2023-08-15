package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDB struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty"`
	Email   string               `bson:"email"`
	Degrees []primitive.ObjectID `bson:"degrees"`
}

type UserStorage struct {
	db *mongo.Database
}

func NewUserStorage(db *mongo.Database) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) FindUserByEmail(email string) (*UserDB, error) {
	var user UserDB
	filter := bson.M{"email": email}
	err := s.db.Collection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, ErrDatabaseFetchUser
	}
	return &user, nil
}

func (s *UserStorage) CreateNewUser(email string) (string, error) {
	user := UserDB{
		Email:   email,
		Degrees: []primitive.ObjectID{},
	}
	res, err := s.db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return "", ErrDatabaseCreateUser
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}
