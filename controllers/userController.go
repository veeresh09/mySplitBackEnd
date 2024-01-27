package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"mySplitBackEnd/config"
	"mySplitBackEnd/models"
	"net/http"
	"time"
)

// ExampleAPIHandler handles the /api/example endpoint
func ExampleAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "This is an example API endpoint from controllers"}`))
}

func CreateUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	var user models.User

	// Decode the request body into the user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exists, err := userExists(collection, user.Email, user.MobileNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User with the given email or mobile number already exists", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert the user into the collection
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

func userExists(collection *mongo.Collection, email, mobileNumber string) (bool, error) {
	var result models.User
	filter := bson.M{
		"$or": []bson.M{
			{"email": email},
			{"mobileNumber": mobileNumber},
		},
	}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	// If an error other than ErrNoDocuments occurs, return the error
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return false, err
	}

	// If no documents are found, the user does not exist
	return !errors.Is(err, mongo.ErrNoDocuments), nil
}

// SignIn handles user authentication and returns a JWT.
func SignIn(w http.ResponseWriter, r *http.Request, collection *mongo.Collection, groupCollection *mongo.Collection, expenseCollection *mongo.Collection) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find user by email
	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &models.Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtKey) // jwtKey is your secret key

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var groups []models.Group
	groupFilter := bson.M{"users": bson.M{"$in": []interface{}{user.ID}}}
	groupCursor, err := groupCollection.Find(context.TODO(), groupFilter)
	if err != nil {
		http.Error(w, "Error fetching groups: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(groupCursor *mongo.Cursor, ctx context.Context) {
		err := groupCursor.Close(ctx)
		if err != nil {

		}
	}(groupCursor, context.TODO())
	for groupCursor.Next(context.TODO()) {
		var group models.Group
		if err := groupCursor.Decode(&group); err != nil {
			http.Error(w, "Error decoding group: "+err.Error(), http.StatusInternalServerError)
			return
		}
		groups = append(groups, group)
	}

	// Find expenses created by the user
	var expenses []models.Expense
	expenseFilter := bson.M{"createdBy": user.ID}
	expenseCursor, err := expenseCollection.Find(context.TODO(), expenseFilter)
	if err != nil {
		http.Error(w, "Error fetching expenses: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(expenseCursor *mongo.Cursor, ctx context.Context) {
		err := expenseCursor.Close(ctx)
		if err != nil {

		}
	}(expenseCursor, context.TODO())
	for expenseCursor.Next(context.TODO()) {
		var expense models.Expense
		if err := expenseCursor.Decode(&expense); err != nil {
			http.Error(w, "Error decoding expense: "+err.Error(), http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, expense)
	}

	// Find other users in these groups
	userIDsSet := make(map[primitive.ObjectID]struct{})
	for _, group := range groups {
		for _, userID := range group.Users {
			userIDsSet[userID] = struct{}{}
		}
	}

	var usersInGroups []models.User
	for userID := range userIDsSet {
		if userID == user.ID {
			continue // Skip the current user
		}
		var user models.User
		if err := collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user); err != nil {
			// Handle error, note that this might fail if the user doesn't exist
			continue
		}
		usersInGroups = append(usersInGroups, user)
	}

	// Return the token and additional data
	response := struct {
		Token         string           `json:"token"`
		Groups        []models.Group   `json:"groups"`
		Expenses      []models.Expense `json:"expenses"`
		UsersInGroups []models.User    `json:"usersInGroups"`
	}{
		Token:         tokenString,
		Groups:        groups,
		Expenses:      expenses,
		UsersInGroups: usersInGroups,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
