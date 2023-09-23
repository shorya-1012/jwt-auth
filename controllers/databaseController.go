package controllers

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "jwt-auth"
const collectionName = "users"

var collection *mongo.Collection

func ConnectToDB(){
    if err := godotenv.Load() ; err != nil {
        log.Fatal(err)
    }

    databaseURI := os.Getenv("MONGODB_URI")
    if databaseURI == ""{
        log.Fatal("Database uri not found in env file")
    }

    clientOptions := options.Client().ApplyURI(databaseURI)

    client , err := mongo.Connect(context.TODO() , clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    collection = client.Database(dbName).Collection(collectionName)
}
