package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Group represents a group of users
type Group struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty"`
	Name    string               `bson:"name"`
	Users   []primitive.ObjectID `bson:"users"`   // Array of User IDs
	Creator primitive.ObjectID   `bson:"creator"` // ID of the user who created the group
}
