package controller

import (
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(usecase usecase.UserUsecase) UserController {
	return UserController{userUsecase: usecase}
}

func isValidPassword(password string) bool {
	const passwordRegex = `^[A-Za-z\d@$!%*?&]{8,}$`
	re := regexp.MustCompile(passwordRegex)
	if !re.MatchString(password) {
		return false
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	var rawData map[string]interface{}
	err := ctx.ShouldBindJSON(&rawData)
	if err != nil {
		response := model.Response{
			Message: "Invalid JSON",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if name, ok := rawData["name"].(string); !ok || name == "" {
		response := model.Response{
			Message: "Invalid name, must be a non-empty string",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	email, ok := rawData["email"].(string)
	if !ok || email == "" {
		response := model.Response{
			Message: "Invalid email, must be a non-empty string",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		response := model.Response{
			Message: "Email is badly formatted",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	password, ok := rawData["password"].(string)
	if !ok || password == "" {
		response := model.Response{
			Message: "Invalid password, must be a non-empty string",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if !isValidPassword(password) {
		response := model.Response{
			Message: "Password must have, at least one character lowercase, uppercase, number and a symbol.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user := model.User{
		Name:     rawData["name"].(string),
		Email:    rawData["email"].(string),
		Password: rawData["password"].(string),
	}

	insertedUser, err := u.userUsecase.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			response := model.Response{
				Message: "Email already in use",
			}
			ctx.JSON(http.StatusConflict, response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, insertedUser)
}
