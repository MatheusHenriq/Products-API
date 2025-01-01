package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
	"strings"
)

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{connection: connection}
}

func (ur *UserRepository) CreateUser(user model.User) (int, error) {
	var id int
	query, err := ur.connection.Prepare("INSERT INTO users" +
		"(name,email,password)" +
		" VALUES ($1,$2,$3) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		fmt.Println(strings.Split(err.Error(), ":"))
		return 0, err
	}
	query.Close()
	return id, nil
}
