package handler

import (
	schema_handlers "blogApp/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SingnUp(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser schema_handlers.User

		err := json.NewDecoder(r.Body).Decode(&newUser)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		var existingUser schema_handlers.User
		err = db.Collection("users").FindOne(context.Background(), bson.M{"email": newUser.Email}).Decode(&existingUser)
		if err != mongo.ErrNoDocuments{
			if err != nil {
				http.Error(w, "Error checking for existing user", http.StatusInternalServerError)
				return
			}
			http.Error(w, "User already exists", http.StatusConflict)
			return 
		}

		
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if(err != nil){
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		newUser.Password = string(hashedPassword);

		insertUser, err := db.Collection("users").InsertOne(context.Background(), newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(newUser)
		fmt.Println(insertUser)
	}
}

func GetUsers(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		collectionName := params["users"]

		collection := db.Collection(collectionName)
		var users []schema_handlers.User

		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var user schema_handlers.User

			if err := cursor.Decode(&user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
