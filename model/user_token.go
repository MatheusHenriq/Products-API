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

const JWT_SECRET_KEY = "JWT_SECRET_KEY"
const UserUUIDKey = "userUUID"

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

func RefreshToken(c *gin.Context) (string, string, error) {

	secret := os.Getenv(JWT_SECRET_KEY)
	claims, err := getClaims(c)
	if err != nil {
		return "", "", err
	}
	currentToken := getCurrentToken(c)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	rt, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return currentToken, rt, nil
}

func getUuid(c *gin.Context) string {
	tokenString := getCurrentToken(c)
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return ""
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	uuid, ok := claims["uuid"].(string)
	if !ok {
		return ""
	}

	return uuid

}

func getClaims(c *gin.Context) (jwt.Claims, error) {
	tokenString := getCurrentToken(c)
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return token.Claims, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return token.Claims, err
	}
	return claims, nil

}
func getCurrentToken(c *gin.Context) string {
	tokenValue := removeBearerPrefix(c.Request.Header.Get("Authorization"))
	fmt.Println(tokenValue)
	return tokenValue
}
func VerifyTokenMiddleware(c *gin.Context) {
	secret := os.Getenv(JWT_SECRET_KEY)
	tokenValue := getCurrentToken(c)
	token, err := jwt.Parse(removeBearerPrefix(tokenValue), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secret), nil
		}
		return nil, fmt.Errorf("invalid token")
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
	c.Set(UserUUIDKey, getUuid(c))
	c.Next()
}

func removeBearerPrefix(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix("Bearer ", token)
	}
	return token
}
