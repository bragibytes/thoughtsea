package main

import (
	"context"
	"dedpidgon/thoughtsea/controllers"
	"dedpidgon/thoughtsea/models"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PORT string = ":9000"
)

func main() {

	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server := gin.Default()

	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_STRING"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	server.Use(cors.New(config))

	models.Init(client)
	controllers.Core{Engine: server}.Init()

	// running server
	fmt.Println("Server running on port ", PORT)
	server.Run(PORT)

}
