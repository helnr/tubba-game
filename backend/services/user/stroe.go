package user

import (
	"context"
	"fmt"

	"github.com/helnr/tubba-game/backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type store struct {
	db *mongo.Database
}


func NewUserStore(db *mongo.Database) *store {
	return &store{
		db: db,
	}
}


func (s *store) GetUserByID(id primitive.ObjectID) (*types.User, error) {
	result := s.db.Collection("users").FindOne(context.TODO(), bson.M{"_id": id})

	if result.Err() != nil {
		return nil, result.Err()
	}

	var user types.User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *store) GetUserByEmail(email string) (*types.User, error) {

	cursor, err := s.db.Collection("users").Find(context.TODO(), bson.M{"email": email})
	if err != nil {
		return nil, err
	}

	var user types.User
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}


	return nil, fmt.Errorf("User not found")
}

func (s *store) SaveUser(user *types.User) error {
	result, err := s.db.Collection("users").InsertOne(
		context.TODO(),
		user,
	)

	user.ID = result.InsertedID.(primitive.ObjectID)	

	if err != nil {
		return err
	}

	return nil
}	

func (s *store) UpdateUser(user *types.User) error {
	_, err := s.db.Collection("users").UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)

	if err != nil {
		return err
	}

	return nil
}