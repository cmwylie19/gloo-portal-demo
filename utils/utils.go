package utils

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func ConnectDB() *mongo.Collection {
	config := GetConfiguration()
	clientOptions := options.Client().ApplyURI(config.ConnectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("taskdb").Collection("tasks")
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

type Configuration struct {
	Port             string
	ConnectionString string
}

func GetConfiguration() Configuration {
	// err := godotenv.Load("./.env")

	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	return Configuration{
		os.Getenv("PORT"),
		os.Getenv("CONNECTION_STRING"),
	}
}
