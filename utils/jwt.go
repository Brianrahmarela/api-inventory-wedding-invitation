package utils

import (
	"api-go-test/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// sekarang menerima role
func GenerateToken(userId uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(config.GetJwtExpirationDuration()).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.GetJwtSecret())
}

// mengembalikan userId dan role
func ValidateToken(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJwtSecret(), nil
	})
	if err != nil {
		return 0, "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := uint(claims["user_id"].(float64))
		role := claims["role"].(string)
		return userId, role, nil
	}
	return 0, "", fmt.Errorf("invalid token")
}
