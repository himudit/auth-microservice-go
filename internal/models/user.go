package models

import (
	"authService/config"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"-"`
	Role         string             `bson:"role" json:"role"`
	CreatedAt    int64              `bson:"created_at" json:"created_at"`
	TokenVersion int                `bson:"token_version" json:"token_version"` // Added field
}

var userCollection *mongo.Collection

func InitCollections() {
	userCollection = config.MongoDB.Collection("users")
}

// Check if a user exists by email
func IsEmailExists(email string) (bool, error) {
	var user User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// InsertUser inserts a new user into the database
func InsertUser(user *User) error {
	_, err := userCollection.InsertOne(context.TODO(), user)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id string) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user User
	err = userCollection.FindOne(
		context.TODO(),
		bson.M{"_id": objectID},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func IncrementTokenVersion(userID string) error {
	// Convert string to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Increment tokenVersion
	update := bson.M{
		"$inc": bson.M{"token_version": 1},
	}

	_, err = userCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		update,
	)

	return err
}
