package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/shorya-1012/jwt-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func RegisterController(c *fiber.Ctx) error {
	var payload models.User

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(422).JSON(fiber.Map{
			"error": "Not processable",
		})
	}

	if payload.Username == "" || payload.Password == "" {
		return c.Status(422).JSON(fiber.Map{
			"error": "Required Params Not Provided",
		})
	}

	if len(payload.Password) < 6 {
		return c.Status(422).JSON(fiber.Map{
			"error": "Password sould contain atleast 6 characters",
		})
	}

	if len(payload.Username) < 3 {
		return c.Status(422).JSON(fiber.Map{
			"error": "Username sould contain atleast 3 characters",
		})
	}

	//check if user with username already exists

	check, err := collection.CountDocuments(context.Background(), bson.D{{"username", payload.Username}})

	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error occured while querign database",
		})
	}

	if check >= 1 {
		return c.Status(403).JSON(fiber.Map{
			"error": "Username already taken",
		})
	}

	// encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error while hasing password",
		})
	}

	payload.Password = string(hashedPassword)

	// insert new user into db
	newUser, err := collection.InsertOne(context.Background(), payload)

	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error occured while creating record",
		})
	}

	return c.JSON(newUser)
}

func LoginController(c *fiber.Ctx) error {

	accessTokenSecretKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if accessTokenSecretKey == "" {
		log.Fatal("Database uri not found in env file")
	}

	refreshTokenSecretKey := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	if refreshTokenSecretKey == "" {
		log.Fatal("Database uri not found in env file")
	}

	var payload models.User
	var user models.User

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(422).JSON(fiber.Map{
			"error": "provided",
		})
	}

	if payload.Username == "" || payload.Password == "" {
		return c.Status(422).JSON(fiber.Map{
			"error": "Required parameteres not provided",
		})
	}

	filter := bson.D{{"username", payload.Username}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		fmt.Println(err)
		return c.Status(404).JSON(fiber.Map{
			"error": "User does not exist",
		})
	}

	matchPassword := checkPassword(user.Password, payload.Password)

	if !matchPassword {
		return c.Status(403).JSON(fiber.Map{
			"error": "Incorrect Password",
		})
	}

	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.UserId.Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	accessToken, err := generateAccessToken(&user, &accessTokenSecretKey)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error while createing token",
		})
	}

	refreshToken, err := refreshTokenClaims.SignedString([]byte(refreshTokenSecretKey))
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error whlie creating token",
		})
	}

	userUpdateFilter := bson.D{{"_id", user.UserId}}
	update := bson.D{{"$set", bson.D{{"refreshtoken", refreshToken}}}}

	updatedUser, err := collection.UpdateOne(context.TODO(), userUpdateFilter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error while updating document",
		})
	}

	fmt.Println(updatedUser)

	// store refresh token in cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 15),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(accessToken)
}

func LogoutController(c *fiber.Ctx)error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
        return c.SendStatus(204)
	}

	userUpdateFilter := bson.D{{"refreshtoken", cookie}}
    update := bson.D{{"$set", bson.D{{"refreshtoken", ""}}}}

	updatedUser, err := collection.UpdateOne(context.TODO(), userUpdateFilter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error while updating document",
		})
	}

    fmt.Println(updatedUser)
    newCookie := fiber.Cookie{
        Name: "jwt",
        Value: "",
        Expires: time.Now().Add(-time.Hour),
        HTTPOnly: true,
    }

	c.Cookie(&newCookie)


    return c.JSON(fiber.Map{
        "message" : "logout successful",
    })
}

func generateAccessToken(user *models.User, accessTokenSecretKey *string) (string, error) {
    id := user.UserId.Hex()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
	})

	return claims.SignedString([]byte(*accessTokenSecretKey))
}

func checkPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
