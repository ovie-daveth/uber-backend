package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"uber-backend/helpers"
	"uber-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var user *models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Bind User"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.ID = primitive.NewObjectID()
		hashedPasswordStr := string(hashedPassword)
		user.Password = &hashedPasswordStr

		_, _ = userCollection.InsertOne(ctx, user)

		jwtToken, err := helpers.GenerateJWT(*user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		userResponse := models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			ImageUrl:  user.ImageUrl,
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "Way to go, you're in",
				"data":    userResponse,
				"jwt":     jwtToken,
			})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var auth models.LoginRequest // Not using pointer here since it's not needed
		var user models.User

		// Bind the JSON input
		if err := c.BindJSON(&auth); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind loginRequest"})
			return // Add return to stop further execution if there is an error
		}

		fmt.Printf("Binded object: %+v\n", auth)

		// Ensure Email and Password are not nil
		if auth.Email == nil || auth.Password == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Password are required"})
			return
		}

		// Find the user by email
		err := userCollection.FindOne(ctx, bson.M{"email": *auth.Email}).Decode(&user)
		fmt.Printf("Binded object: %+v\n", *auth.Email)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching user details"})
			return
		}

		// Ensure user.Password is not nil
		if user.Password == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User password not set"})
			return
		}

		// Compare hashed password with the provided password
		err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*auth.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		jwtToken, err := helpers.GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		userResponse := models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			ImageUrl:  user.ImageUrl,
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{
			"message": "Hurray, it is you",
			"data":    userResponse,
			"jwt":     jwtToken,
		})
	}
}

func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the JWT from the Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("bearer %s", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			return
		}

		// Token should be in the form "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Printf("tokenString %s", tokenString)
		if tokenString == authHeader {
			// If token doesn't contain "Bearer"
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token format is invalid"})
			return
		}

		// Parse and validate the token
		claims, err := helpers.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract the email from the claims
		email := claims["email"].(string)
		fmt.Printf("bearer %s", email)

		// Use the email to fetch the user details from MongoDB
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
			return
		}

		userResponse := models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			ImageUrl:  user.ImageUrl,
		}

		// Return user profile details
		c.JSON(http.StatusOK, gin.H{
			"message": "Profile details",
			"data":    userResponse,
		})
	}
}

func UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			return
		}

		// Token should be in the form "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token format is invalid"})
			return
		}

		// Parse and validate the token
		claims, err := helpers.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract the email from the claims
		email := claims["email"].(string)

		// Use the email to fetch the user details from MongoDB
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var userRequest models.UserResponse
		if err = c.BindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind user details"})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"firstname": userRequest.FirstName,
				"lastname":  userRequest.LastName,
				"username":  userRequest.Username,
				"email":     userRequest.Email,
				"imageurl":  userRequest.ImageUrl,
			},
		}

		_, err = userCollection.UpdateOne(ctx, bson.M{"email": email}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user details"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully"})
	}
}
