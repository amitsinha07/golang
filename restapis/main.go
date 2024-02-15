package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello, world!")
}

func createUserHandler(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Parse the request body into a new User instance
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(newUser)
	// Insert the new user into the collection
	insertResult, err := collection.InsertOne(context.Background(), newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the ID of the newly inserted user back in the response
	json.NewEncoder(w).Encode(newUser)
	fmt.Println(insertResult)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var users []*User
	if err = cursor.All(context.Background(), &users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Parse the request body into a new User instance
	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure we have an ID to update
	if updatedUser.ID == primitive.NilObjectID {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Create a filter to select the user by ID
	filter := bson.M{"_id": updatedUser.ID}

	// Update the user document in the collection
	updateResult, err := collection.UpdateOne(context.Background(), filter, bson.M{"$set": updatedUser})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the update was successful
	if updateResult.MatchedCount ==  0 {
		http.Error(w, "No user found with the given ID", http.StatusNotFound)
		return
	}

	// Send a success response
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Successfully updated %d user(s)", updateResult.ModifiedCount),
	})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Extract the user ID from the query parameters
	queryParams := r.URL.Query()
	idStr := queryParams.Get("id")

	// Validate the ID parameter
	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	// Convert the ID string to a primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Create a filter to select the user by ID
	filter := bson.M{"_id": id}

	// Delete the user document from the collection
	deleteResult, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the delete was successful
	if deleteResult.DeletedCount ==  0 {
		http.Error(w, "No user found with the given ID", http.StatusNotFound)
		return
	}

	// Send a success response
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "User successfully deleted",
	})
}



func getUsersCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("users")
}

func main() {
	client, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("golangDB")
	collection := getUsersCollection(db)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsersHandler(w, r, collection)
		case http.MethodPost:
			createUserHandler(w, r, collection)
		case http.MethodPut:
			updateUserHandler(w, r, collection)
		case http.MethodDelete:
			deleteUserHandler(w, r, collection)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the server on port  8080
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// MongoDB connection
func connectDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
