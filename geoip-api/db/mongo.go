package db

import (
	"context"
	"time"

	"github.com/konstfish/angler/geoip-api/configs"
	"github.com/konstfish/angler/geoip-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	client = ConnectMongo()
}

func ConnectMongo() *mongo.Client {
	uri := configs.GetConfigVar("MONGODB_URI")

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

// returns false if the address doesn't exist or the last check was more than 2 weeks ago
func CheckAddress(address string) (bool, error) {
	collection := GetCollection("angler", "geoip")

	var result models.GeoIP
	err := collection.FindOne(context.TODO(), bson.M{"_id": address}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	twoWeeksAgo := time.Now().Add(-14 * 24 * time.Hour)
	if result.AddressAge.Before(twoWeeksAgo) {
		return false, nil
	}

	return true, nil
}
