package db

import (
	"context"

	"github.com/konstfish/angler/shared/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var client *mongo.Client

func InitMongo() {
	client = ConnectMongo()
}

func ConnectMongo() *mongo.Client {
	uri := configs.GetConfigVar("MONGODB_URI")

	opts := options.Client().ApplyURI(uri)
	opts.Monitor = otelmongo.NewMonitor()

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	return client
}

func DisconnectMongo() {
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}
