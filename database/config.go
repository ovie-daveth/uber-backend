package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// init() runs automatically during package initialization
func init() {
	Client = DbInstance()
}

func DbInstance() *mongo.Client {

	if Client != nil {
		return Client // Return existing client if already initialized
	}

	mongoDb := os.Getenv("mongodbUrl")
	fmt.Println("Connecting to MongoDB...")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	var err error

	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoDb))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return Client

}

//create collections

func OpenCollection(collectionName string) *mongo.Collection {
	var collection *mongo.Collection = Client.Database("Uber").Collection(collectionName)

	return collection
}
