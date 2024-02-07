package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Created  float64            `bson:"created" json:"created"`
	Settings DomainSettings     `bson:"settings" json:"settings"`
}

type DomainSettings struct {
	EnabledUntil float64 `bson:"enabled_until" json:"enabled_until"`
	Enabled      bool    `bson:"enabled" json:"enabled"`
}

func (domain *Domain) Serialize() string {
	domainJSON, err := json.Marshal(domain)
	if err != nil {
		return ""
	}

	return string(domainJSON)
}

func (domain *Domain) Deserialize(domainJSON string) error {
	err := json.Unmarshal([]byte(domainJSON), &domain)
	return err
}
