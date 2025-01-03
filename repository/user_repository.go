package repository

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"go-api/model"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserRepository struct {
	connection *sql.DB
}

var (
	JWT_SECRET_KEY = "JWT_SECRET_KEY"
)

func encryptPassword(password string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(password))
	password = hex.EncodeToString(hash.Sum(nil))
	return password
}

func generateToken(user model.User) (string, error) {

	secret := os.Getenv(JWT_SECRET_KEY)
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"isAdmin": user.IsAdmin,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {

		return "", err
	}
	return tokenString, nil
}

func removeBearerPrefix(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix("Bearer ", token)
	}
	return token
}

// I prefer to return error, because if error == nil it means that token is correct, if occurs an error, it will be an invalid token
func VerifyToken(tokenValue string) error {
	secret := os.Getenv(JWT_SECRET_KEY)
	token, err := jwt.Parse(removeBearerPrefix(tokenValue), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secret), nil
		}
		return nil, fmt.Errorf("Invalid token")
	})

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return err
	}
	return nil
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{connection: connection}
}

func (ur *UserRepository) CreateUser(user model.User) (int, error) {
	var id int

	var password = encryptPassword(user.Password)
	query, err := ur.connection.Prepare("INSERT INTO users" +
		"(name,email,password)" +
		" VALUES ($1,$2,$3) RETURNING id")
	if err != nil {
		return 0, err
	}

	err = query.QueryRow(user.Name, user.Email, password).Scan(&id)
	if err != nil {
		fmt.Println(strings.Split(err.Error(), ""))
		return 0, err
	}
	query.Close()
	return id, nil
}

func (ur *UserRepository) DeleteUser(user model.User) (*model.User, error) {
	query, err := ur.connection.Prepare("DELETE FROM users" +
		" WHERE email = $1 and password = $2 RETURNING id")

	var userData model.User
	password := encryptPassword(user.Password)
	if err != nil {
		return nil, err
	}

	err = query.QueryRow(user.Email, password).Scan(&userData.ID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	query.Close()
	return &userData, nil
}

func (ur *UserRepository) LogIn(user model.User) (*model.User, string, error) {
	query, err := ur.connection.Prepare("select id,name,email,is_admin from users where email = $1 and password = $2")
	if err != nil {
		return nil, "", err
	}
	var userData model.User
	var password = encryptPassword(user.Password)
	err = query.QueryRow(user.Email, password).Scan(&userData.ID, &userData.Name, &userData.Email, &userData.IsAdmin)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		return nil, "", err
	}
	query.Close()
	token, err := generateToken(userData)
	if err != nil {
		return nil, "", err
	}
	return &userData, token, nil
}
