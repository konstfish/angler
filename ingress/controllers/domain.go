package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/konstfish/angler/ingress/models"
	"github.com/konstfish/angler/shared/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var domainCollection *mongo.Collection
var domainTTL time.Duration = time.Hour * 6

func InitCollections() {
	domainCollection = db.GetCollection("angler", "domains")
	eventCollection = db.GetCollection("angler", "events")
	sessionCollection = db.GetCollection("angler", "sessions")
}

func ValidateDomain(ctx context.Context, domainName string) (bool, error) {
	domain, err := getCacheDomain(ctx, domainName)
	if err != nil {
		return false, err
	}

	if domain.Settings.EnabledUntil < float64(time.Now().UnixMilli()) {
		return false, nil
	}

	// todo: amount of sessions this month. but also need db so this is fine for now

	return domain.Settings.Enabled, nil
}

func getCacheDomain(ctx context.Context, domainName string) (models.Domain, error) {
	domainJSON, err := db.Redis.Client.Get(ctx, fmt.Sprintf("dm-%s", domainName)).Result()
	if err != nil {
		log.Println("cache miss")
		domain, err := getDomain(ctx, domainName)

		return domain, err
	}

	_, err = db.Redis.Client.Expire(ctx, fmt.Sprintf("dm-%s", domainName), domainTTL).Result()

	var domain models.Domain
	domain.Deserialize(domainJSON)
	if err != nil {
		return domain, err
	}

	return domain, nil
}

func getDomain(ctx context.Context, domainName string) (models.Domain, error) {
	var domain models.Domain

	err := domainCollection.FindOne(ctx, bson.M{"name": domainName}).Decode(&domain)
	if err != nil {
		return domain, err
	}

	go writeCacheDomain(ctx, domain)
	return domain, nil
}

func writeCacheDomain(ctx context.Context, domain models.Domain) {
	db.Redis.Client.Set(ctx, fmt.Sprintf("dm-%s", domain.Name), domain.Serialize(), domainTTL)
}
