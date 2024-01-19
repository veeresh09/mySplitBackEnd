package controllers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mySplitBackEnd/models"
	"net/http"
)

// CreateGroup handles the creation of a new group.
func CreateGroup(w http.ResponseWriter, r *http.Request, userCollection *mongo.Collection, groupCollection *mongo.Collection) {
	var request struct {
		Name    string   `json:"name"`
		Emails  []string `json:"emails"`
		Creator string   `json:"creator"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if creator is specified
	if request.Creator == "" {
		http.Error(w, "Creator must be specified", http.StatusBadRequest)
		return
	}

	// Initialize group with unique members and add creator if not present
	uniqueMembers := make(map[string]primitive.ObjectID)
	for _, email := range request.Emails {
		var user models.User
		if err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user); err == nil {
			uniqueMembers[email] = user.ID
		}
	}

	// Check if creator is in the list of emails, add if not
	creatorID, err := getCreatorID(userCollection, request.Creator)
	if err != nil {
		http.Error(w, "Invalid creator", http.StatusBadRequest)
		return
	}
	uniqueMembers[request.Creator] = creatorID

	// Convert map keys to slice of ObjectIDs
	users := make([]primitive.ObjectID, 0, len(uniqueMembers))
	for _, id := range uniqueMembers {
		users = append(users, id)
	}

	// Create and insert the group
	group := models.Group{
		ID:      primitive.NewObjectID(),
		Name:    request.Name,
		Users:   users,
		Creator: creatorID,
	}
	_, err = groupCollection.InsertOne(context.TODO(), group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created group
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		return
	}
}

// getCreatorID retrieves the ID of the creator by their email.
func getCreatorID(collection *mongo.Collection, email string) (primitive.ObjectID, error) {
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return user.ID, nil
}
