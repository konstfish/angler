package db

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var db *geoip2.Reader

func init() {
	var err error

	db, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatalf("Failed to create GeoIP database reader: %v", err)
	}
}

func GetAddressData(address string) (*geoip2.City, error) {
	ip := net.ParseIP(address)

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	return record, err
}
