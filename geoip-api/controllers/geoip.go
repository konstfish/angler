package controllers

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	geodb "github.com/konstfish/angler/geoip-api/db"
	"github.com/konstfish/angler/geoip-api/models"
	"github.com/konstfish/angler/shared/db"
	"github.com/konstfish/angler/shared/monitoring"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/trace"
)

var geoIpCollection *mongo.Collection

func InitCollections() {
	geoIpCollection = db.GetCollection("angler", "geoip")
}

func isValidAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

func GetIpInfo(ctx context.Context, address string) (models.GeoIP, error) {
	var span trace.Span
	ctx, span = monitoring.Tracer.Start(ctx, "GetIpInfo")
	defer span.End()

	if !isValidAddress(address) {
		return models.GeoIP{}, errors.New("Invalid IP address")
	}

	record, err := geodb.GetAddressData(address)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(record)

	var region string
	if len(record.Subdivisions) > 0 {
		if name, ok := record.Subdivisions[0].Names["en"]; ok {
			region = name
		}
	}

	geoip := models.GeoIP{
		Address:     address,
		AddressAge:  time.Now(),
		CountryCode: record.Country.IsoCode,
		CountryName: record.Country.Names["en"],
		RegionName:  region,
		CityName:    record.City.Names["en"],
	}

	return geoip, err
}

func ProcessAddress(ctx context.Context, address string) {
	exists, err := CheckAddress(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		return
	}

	geoip, err := GetIpInfo(ctx, address)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": geoip.Address}
	opts := options.Replace().SetUpsert(true)
	result, err := geoIpCollection.ReplaceOne(ctx, filter, geoip, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Added IP:", geoip.Address, result)
}

// TODO implement cache check
// returns false if the address doesn't exist or the last check was more than 2 weeks ago
func CheckAddress(ctx context.Context, address string) (bool, error) {
	var result models.GeoIP
	err := geoIpCollection.FindOne(ctx, bson.M{"_id": address}).Decode(&result)
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
