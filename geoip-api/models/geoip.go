package models

import "time"

type GeoIP struct {
	Address     string    `bson:"address" json:"ip"`
	AddressAge  time.Time `bson:"address_age" json:"address_age"`
	CountryCode string    `bson:"country_code" json:"country_code"`
	CountryName string    `bson:"country_name" json:"country_name"`
	RegionName  string    `bson:"region_name" json:"region_name"`
	CityName    string    `bson:"city_name" json:"city_name"`
}
