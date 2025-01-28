package game

import (
	"context"
	"log"

	"github.com/helnr/tubba-game/backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mongoStore struct {
	db *mongo.Database
	collection *mongo.Collection
}

func NewMongoGameStore(db *mongo.Database) *mongoStore {
	collection := db.Collection("games")

	indexOptions := options.Index().SetExpireAfterSeconds(60 * 60 * 24)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"createdAt": 1},
		Options: indexOptions,
	}

	res, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Println("Error creating index:", err)
	}
	log.Println("Index created:", res)

	return &mongoStore{
		db: db,
		collection: collection,
	}
}

func (g *mongoStore) GetGameByID(id primitive.ObjectID) (*types.Game, error) {
	result := g.collection.FindOne(context.TODO(), bson.M{"_id": id})
	
	if result.Err() != nil {
		return nil, result.Err()
	}

	var game types.Game
	err := result.Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
	
}

func (g *mongoStore) SaveGame(game *types.Game) error {
	result, err := g.collection.InsertOne(context.TODO(), game)
	if err != nil {
		return err
	}
	game.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}


func (g *mongoStore) UpdateGame(game *types.Game) error {
	filter := bson.M{"_id": game.ID}
	update := bson.M{"$set": game}
	_, err := g.collection.UpdateOne(context.TODO(), filter, update)
	return err
}