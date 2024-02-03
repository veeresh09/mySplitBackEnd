package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mySplitBackEnd/controllers"
	"mySplitBackEnd/db"
	"net/http"
)

func main() {
	client := db.Connect()
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println(err)
		}
	}(client, context.TODO())
	usersCollection := db.GetUsersCollection(client)
	groupCollection := db.GetGroupsCollection(client)
	expenseCollection := db.GetExpenseCollection(client)
	r := mux.NewRouter()
	r.HandleFunc("/api/example", controllers.ExampleAPIHandler)

	r.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateUser(w, r, usersCollection)
	}).Methods("POST")

	r.HandleFunc("/api/signin", func(w http.ResponseWriter, r *http.Request) {
		controllers.SignIn(w, r, usersCollection, groupCollection, expenseCollection)
	}).Methods("POST")

	r.HandleFunc("/api/user/email", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUserByEmail(w, r, usersCollection)
	}).Methods("GET")

	r.HandleFunc("/api/user/phoneNumber", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUserByPhoneNumber(w, r, usersCollection)
	}).Methods("GET")

	r.HandleFunc("/api/groups", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateGroup(w, r, usersCollection, groupCollection) // Assuming groupCollection is defined
	}).Methods("POST")

	r.HandleFunc("/api/expenses", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateExpense(w, r, expenseCollection)
	}).Methods("POST")

	r.HandleFunc("/api/expenses/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetExpense(w, r, expenseCollection)
	}).Methods("GET")

	r.HandleFunc("/api/expenses/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateExpense(w, r, expenseCollection)
	}).Methods("PUT")

	r.HandleFunc("/api/expenses/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteExpense(w, r, expenseCollection)
	}).Methods("DELETE")

	r.HandleFunc("/api/groups/{groupId}/expenses", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetExpensesByGroup(w, r, expenseCollection)
	}).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
