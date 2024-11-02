package controller

import (
	"uber-backend/database"

	"go.mongodb.org/mongo-driver/mongo"
)

// Centralizing collection references to avoid duplication
var userCollection *mongo.Collection = database.OpenCollection("user")
