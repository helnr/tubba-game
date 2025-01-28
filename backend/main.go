package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/helnr/tubba-game/backend/config"
	"github.com/helnr/tubba-game/backend/db"
	"github.com/helnr/tubba-game/backend/services/game"
	"github.com/helnr/tubba-game/backend/services/user"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client := db.NewMongoClient(
		options.
		Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.Env.DBHost, config.Env.DBPort)))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := client.Ping(ctx, nil)
	if err != nil {
		log.Println("Unable to connect to database:")
		log.Fatal(err)
	}
	log.Println("Connected to database")

	database := db.NewMongoDatabase(client, "tubba")

	fiber := fiber.New()
	fiber.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	userStore := user.NewUserStore(database)
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(fiber)

	gameStore := game.NewMongoGameStore(database)
	gameService := game.NewGameHandler(userStore, gameStore)
	gameService.RegisterRoutes(fiber)

	fiber.Listen(fmt.Sprintf(":%s", config.Env.ServerPort))

}