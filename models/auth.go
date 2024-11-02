package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"-" bson:"_id"`
	Username  *string            `json:"username" validate:"required, min=2,max=100"`
	Email     *string            `json:"email" validate:"required"`
	FirstName *string            `json:"firstname"`
	LastName  *string            `json:"lastname"`
	Password  *string            `json:"password" validate:"required"`
	ImageUrl  *string            `json:"imageurl"`
}

type LoginRequest struct {
	Email    *string `json:"email" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  *string            `json:"username"`
	Email     *string            `json:"email"`
	FirstName *string            `json:"firstname"`
	LastName  *string            `json:"lastname"`
	ImageUrl  *string            `json:"imageurl"`
}
