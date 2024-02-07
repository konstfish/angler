package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/backend/models"
	"github.com/konstfish/angler/shared/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var domainCollection *mongo.Collection
var domainBindCollection *mongo.Collection

func init() {
	domainCollection = db.GetCollection("angler", "domains")
	domainBindCollection = db.GetCollection("angler", "domainBind")
}

func PostDomain(c *gin.Context) {
	ctx := c.Request.Context()

	// get domain from request
	var domain models.Domain
	if err := c.BindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	// adjust domain fields
	domain.Created = float64(time.Now().UnixMilli())
	domain.ID = primitive.NewObjectID()
	// enabled for a month
	domain.Settings = models.DomainSettings{
		EnabledUntil: float64(time.Now().AddDate(0, 1, 0).UnixMilli()),
		Enabled:      true,
	}

	// insert domain into database
	result, err := domainCollection.InsertOne(ctx, domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// associate domain with user
	user := GetUserFromContext(c)
	result, err = bindUserToDomain(ctx, domain.ID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func bindUserToDomain(ctx context.Context, domainID primitive.ObjectID, userID primitive.ObjectID) (*mongo.InsertOneResult, error) {
	var domainBind models.DomainBind

	domainBind.ID = primitive.NewObjectID()
	domainBind.DomainID = domainID
	domainBind.UserID = userID

	result, err := domainBindCollection.InsertOne(ctx, domainBind)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetDomain(c *gin.Context) {
	ctx := c.Request.Context()

	domainName := c.Param("domain")

	var domain models.Domain
	err := domainCollection.FindOne(ctx, bson.M{"name": domainName}).Decode(&domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain)
}

// TODO get domains owned by user
func GetDomains(c *gin.Context) {
	ctx := c.Request.Context()

	var domains []models.Domain
	user := GetUserFromContext(c)

	// fix yellow underline: https://stackoverflow.com/questions/54548441/composite-literal-uses-unkeyed-fields

	// Open an aggregation cursor
	cursor, err := domainBindCollection.Aggregate(ctx, bson.A{
		bson.D{{"$match", bson.D{{"user_id", user.ID}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "domains"},
					{"localField", "domain_id"},
					{"foreignField", "_id"},
					{"as", "domain"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"domain",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$domain",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$domain"}}}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// iterate over cursor and append to domains slice
	for cursor.Next(ctx) {
		var domain models.Domain
		if err := cursor.Decode(&domain); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		domains = append(domains, domain)
	}

	c.JSON(http.StatusOK, domains)
}
