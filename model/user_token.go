package model

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_SECRET_KEY = "JWT_SECRET_KEY"
)

func GenerateToken(user User) (string, error) {

	secret := os.Getenv(JWT_SECRET_KEY)
	claims := jwt.MapClaims{
		"email":   user.Email,
		"name":    user.Name,
		"uuid":    user.Uuid,
		"isAdmin": user.IsAdmin,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {

		return "", err
	}
	return tokenString, nil
}

func GetUuid(c *gin.Context) (string, error) {
	tokenString := getCurrentToken(c)
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("Failed to parse token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to cast claims to MapClaims")
	}
	uuid, ok := claims["uuid"].(string)
	if !ok {
		return "", fmt.Errorf("uuid field not found or not a string")
	}
	return uuid, nil

}
func getCurrentToken(c *gin.Context) string {

	tokenValue := removeBearerPrefix(c.Request.Header.Get("Authorization"))
	return tokenValue
}
func VerifyTokenMiddleware(c *gin.Context) {
	secret := os.Getenv(JWT_SECRET_KEY)
	tokenValue := getCurrentToken(c)
	token, err := jwt.Parse(removeBearerPrefix(tokenValue), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secret), nil
		}
		return nil, fmt.Errorf("Invalid token")
	})
	response := Response{
		Message: "Invalid token"}
	if err != nil {

		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}
}

func removeBearerPrefix(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix("Bearer ", token)
	}
	return token
}
