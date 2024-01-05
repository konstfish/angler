package db

import (
	"context"
	"os"

	"github.com/konstfish/angler/ingress/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	configs.LoadDotEnv()

	client = ConnectMongo()
}

func ConnectMongo() *mongo.Client {
	uri := os.Getenv("MONGODB_URI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func DisconnectMongo() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}
