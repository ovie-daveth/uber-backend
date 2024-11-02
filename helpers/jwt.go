package helpers

import (
	"errors"
	"log"
	"os"
	"time"
	"uber-backend/models"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(user models.User) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT secret key is not set in the environment")
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = *user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT secret key is not set in the environment")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	// If parsing or validation failed
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims if the token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not extract claims")
	}

	return claims, nil
}
