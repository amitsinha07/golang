package schema_handlers

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name	string	`json:"name" bson:"name"`
	Password	string `json:"password" bson:"password"`
	Email	string	`json:"email" bson:"email"`
}