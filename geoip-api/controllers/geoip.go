package controllers

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/konstfish/angler/geoip-api/db"
	"github.com/konstfish/angler/geoip-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var geoIpCollection *mongo.Collection

func init() {
	geoIpCollection = db.GetCollection("angler", "geoip")
}

func isValidAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

func GetIpInfo(address string) (models.GeoIP, error) {
	if !isValidAddress(address) {
		return models.GeoIP{}, errors.New("Invalid IP address")
	}

	record, err := db.GetAddressData(address)
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

func PushIpInfo(address string) {
	geoip, err := GetIpInfo(address)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": geoip.Address}
	opts := options.Replace().SetUpsert(true)
	result, err := geoIpCollection.ReplaceOne(context.TODO(), filter, geoip, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Added IP:", geoip.Address, result)
}
