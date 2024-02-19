package handler

import (
	schema_handlers "blogApp/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		if err != mongo.ErrNoDocuments {
			if err != nil {
				http.Error(w, "Error checking for existing user", http.StatusInternalServerError)
				return
			}
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newUser.Password = string(hashedPassword)

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

func CreatePost(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPost schema_handlers.BlogPost
		err := json.NewDecoder(r.Body).Decode(&newPost)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newPost.CreatedAt = time.Now()

		insertPost, err := db.Collection("posts").InsertOne(context.Background(), newPost)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(insertPost)
		fmt.Println(insertPost)

	}
}

func GetPosts(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var posts []schema_handlers.BlogPost
		collection := db.Collection("posts")
		cursor, err := collection.Find(context.Background(), bson.M{})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var post schema_handlers.BlogPost

			if err := cursor.Decode(&post); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func UpdatePost(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		postId := params["postid"]

		objId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
		}

		var updatedPost schema_handlers.BlogPost

		err = json.NewDecoder(r.Body).Decode(&updatedPost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		filter := bson.M{"_id": objId}
		update := bson.M{
			"$set": bson.M{
				"title":       updatedPost.Title,
				"description": updatedPost.Description,
			},
		}

		updatedResult, err := db.Collection("posts").UpdateOne(context.Background(), filter, update)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if updatedResult.MatchedCount == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(updatedResult)
		fmt.Println(updatedResult)
	}
}

func DeletePost(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		postId := params["postid"]

		objId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			http.Error(w, "Invalid post Id", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": objId}

		deleteResult, err := db.Collection("posts").DeleteOne(context.Background(), filter)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if deleteResult.DeletedCount == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(deleteResult)
		fmt.Println(deleteResult)
	}
}
