package controllers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/shorya-1012/jwt-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchUserController(c *fiber.Ctx) error {
	accessTokenSecretKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	var user models.User
    var responseUser models.ResponseUser

	header := c.Request().Header.Peek("Authorization")
	accessToken := strings.Split(string(header), " ")[1]

	if accessToken == "" {
		return c.Status(403).JSON(fiber.Map{
			"error": "accecss token not found",
		})
	}

	verifiedToken, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(accessTokenSecretKey), nil
	})
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Token expired",
		})
	}

	claims := verifiedToken.Claims.(*jwt.StandardClaims)

    id , err := primitive.ObjectIDFromHex(claims.Issuer)
    if err != nil {
        fmt.Println(err)
        return c.SendString("error while converting to object id")
    }

    filter := bson.M{"_id" : id}
	error := collection.FindOne(context.TODO(), filter).Decode(&user)

    if error != nil {
        return c.Status(404).JSON(fiber.Map{
            "error" : "user not found",
        })
    }

    responseUser.UserId = id;
    responseUser.Username = user.Username

	return c.JSON(responseUser)
}

