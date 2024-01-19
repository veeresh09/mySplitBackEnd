package controllers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mySplitBackEnd/models"
	"net/http"
	"time"
)

// CreateExpense handles the creation of a new expense.
func CreateExpense(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	var expense models.Expense

	// Decode the request body into the expense struct
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the ID and timestamps
	expense.ID = primitive.NewObjectID()
	expense.CreatedAt = time.Now()
	expense.ModifiedAt = time.Now()

	// Insert the expense into the collection
	_, err = collection.InsertOne(context.TODO(), expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created expense
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(expense)
	if err != nil {
		return
	}
}

// GetExpense retrieves a single expense by its ID.
func GetExpense(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	idParam := mux.Vars(r)["id"] // Get ID from URL
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var expense models.Expense
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

// UpdateExpense updates an existing expense.
func UpdateExpense(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	idParam := mux.Vars(r)["id"]
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var expense models.Expense
	err = json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expense.ModifiedAt = time.Now()
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": expense})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteExpense deletes an expense.
func DeleteExpense(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	idParam := mux.Vars(r)["id"]
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetExpensesByGroup retrieves all expenses for a specific group.
func GetExpensesByGroup(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Extract the group ID from URL parameters
	groupIDParam := mux.Vars(r)["groupId"]
	groupID, err := primitive.ObjectIDFromHex(groupIDParam)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Find all expenses for the group
	cursor, err := collection.Find(context.TODO(), bson.M{"groupId": groupID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and decode each document
	var expenses []models.Expense
	for cursor.Next(context.TODO()) {
		var expense models.Expense
		err := cursor.Decode(&expense)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, expense)
	}

	// Check for any errors encountered during iteration
	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the list of expenses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}
