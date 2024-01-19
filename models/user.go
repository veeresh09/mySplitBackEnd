// Package models user.go
package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Name         string               `bson:"name"`
	MobileNumber string               `bson:"mobileNumber"`
	Email        string               `bson:"email"`
	Groups       []primitive.ObjectID `bson:"groups"` // Array of Group IDs
	Password     string               `bson:"password"`
}

type Claims struct {
	UserID primitive.ObjectID `json:"userId"`
	jwt.StandardClaims
}
