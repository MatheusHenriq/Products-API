package controller

import (
	"go-api/model"
	"go-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(usecase usecase.UserUsecase) UserController {
	return UserController{userUsecase: usecase}
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
	if email, ok := rawData["email"].(string); !ok || email == "" {
		response := model.Response{
			Message: "Invalid email, must be a non-empty string",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if password, ok := rawData["password"].(string); !ok || password == "" {
		response := model.Response{
			Message: "Invalid password, must be a non-empty string",
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

		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, insertedUser)
}
