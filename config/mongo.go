package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongo() {
	
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	log.Printf(uri)

	if uri == "" {
		log.Fatal("❌ MONGO_URI not found in environment")
	}

	if dbName == "" {
		log.Fatal("❌ MONGO_DB not found in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ MongoDB Connection Failed:", err)
	}

	// Ping to check connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB Ping Failed:", err)
	}

	log.Printf("✅ Connected to MongoDB: %s", dbName)

	MongoClient = client
	MongoDB = client.Database(dbName)
}
