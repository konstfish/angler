package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/konstfish/angler/auth/models"
	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userTTL time.Duration = time.Hour * 24 * 2

func PostRegister(c *gin.Context) {
	ctx := c.Request.Context()

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	// todo i have no clue what bycrypt does exactly so reseach it sometime in the future
	user.HashPassword()
	user.ID = primitive.NewObjectID()

	if err := createUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"registered": true})
}

func createUser(ctx context.Context, user models.User) error {
	_, err := userCollection.InsertOne(ctx, user)
	return err
}

func getUser(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := userCollection.FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	return user, err
}

func PostLogin(c *gin.Context) {
	ctx := c.Request.Context()

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	nuser, err := authUser(ctx, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := createJwt(ctx, nuser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": token})
}

func authUser(ctx context.Context, email string, password string) (models.User, error) {
	user, err := getUser(ctx, email)
	if err != nil {
		return models.User{}, err
	}

	if ok, err := user.CheckPassword(password); !ok {
		return models.User{}, err
	}

	return user, nil
}

func createJwt(ctx context.Context, user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Email,
		"sub":  user.ID.Hex(),
		"exp":  time.Now().Add(userTTL).Unix(),
	})

	tokenString, err := token.SignedString([]byte(configs.GetConfigVar("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	cacheJwt(ctx, user.ID.Hex(), tokenString)

	return tokenString, nil
}

func cacheJwt(ctx context.Context, userId string, token string) {
	db.Redis.Client.Set(ctx, fmt.Sprintf("token-%s", userId), token, userTTL)
}

func existsJwt(ctx context.Context, userId string, token string) bool {
	val, err := db.Redis.Client.Get(ctx, fmt.Sprintf("token-%s", userId)).Result()
	if err != nil {
		return false
	}

	return val == token
}

func GetVerify(c *gin.Context) {
	// get bearer token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}

	// remove bearer prefix
	tokenString = tokenString[7:]

	// get user id from token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.GetConfigVar("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// check if token is valid
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// check if token is in cache
	userId := claims["sub"].(string)
	if !existsJwt(c.Request.Context(), userId, tokenString) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// return user id and email
	c.JSON(http.StatusOK, gin.H{
		"id":   userId,
		"mail": claims["user"],
	})
}
