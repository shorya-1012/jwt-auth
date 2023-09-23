package controllers

import (
	"context"
	"fmt"
	"os"
    "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/shorya-1012/jwt-auth/models"
	"go.mongodb.org/mongo-driver/bson"
)

func RefreshTokenController(c *fiber.Ctx) error {
	refreshTokenSecretKey := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	accessTokenSecretKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	var user models.User

	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.JSON(fiber.Map{
			"error": "User not logged in",
		})
	}

	filter := bson.D{{"refreshtoken", cookie}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		//delete refresh token cookie
		newCookie := fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
		}

		c.Cookie(&newCookie)

		return c.Status(404).JSON(fiber.Map{
			"error": "Not found",
		})
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenSecretKey), nil
	})

	fmt.Println(token)

	if err != nil {
		// logout user
		userUpdateFilter := bson.D{{"refreshtoken", cookie}}
		update := bson.D{{"$set", bson.D{{"refreshtoken", ""}}}}

		updatedUser, err := collection.UpdateOne(context.TODO(), userUpdateFilter, update)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "error while updating document",
			})
		}

		fmt.Println(updatedUser)
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accessToken, err := generateAccessToken(&user, &accessTokenSecretKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(accessToken)
}
