package repository

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"go-api/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserRepository struct {
	connection *sql.DB
}

var (
	JWT_SECRET_KEY = "JWT_SECRET_KEY"
)

func GenerateUserID() string {
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	randomHex := hex.EncodeToString(randomBytes)

	return randomHex[:7]
}

func encryptPassword(password string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(password))
	password = hex.EncodeToString(hash.Sum(nil))
	return password
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{connection: connection}
}

func (ur *UserRepository) CreateUser(user model.User) error {
	uuid := GenerateUserID()
	var id int
	var password = encryptPassword(user.Password)
	query, err := ur.connection.Prepare("INSERT INTO users" +
		"(name, email, password, uuid)" +
		" VALUES ($1,$2,$3,$4) RETURNING id")
	if err != nil {
		return err
	}

	err = query.QueryRow(user.Name, user.Email, password, uuid).Scan(&id)
	if err != nil {
		fmt.Println(strings.Split(err.Error(), ""))
		return err
	}
	query.Close()
	return nil
}

func (ur *UserRepository) DeleteUser(user model.User) (*model.User, error) {
	query, err := ur.connection.Prepare("DELETE FROM users" +
		" WHERE email = $1 and password = $2 RETURNING id")

	var userData model.User
	password := encryptPassword(user.Password)
	if err != nil {
		return nil, err
	}

	err = query.QueryRow(user.Email, password).Scan(&userData.Uuid)
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
	query, err := ur.connection.Prepare("select uuid,name,email,is_admin from users where email = $1 and password = $2")
	if err != nil {
		return nil, "", err
	}
	var userData model.User
	var password = encryptPassword(user.Password)
	err = query.QueryRow(user.Email, password).Scan(&userData.Uuid, &userData.Name, &userData.Email, &userData.IsAdmin)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		return nil, "", err
	}
	query.Close()
	token, err := model.GenerateToken(userData)
	if err != nil {
		return nil, "", err
	}
	return &userData, token, nil
}

func (ur *UserRepository) RefreshToken(c *gin.Context) (string, string, error) {
	t, rt, err := model.RefreshToken(c)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	return t, rt, nil
}
