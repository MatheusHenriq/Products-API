package main

import (
	"go-api/controller"
	"go-api/db"
	"go-api/model"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	//camada de repository
	ProductRepository := repository.NewProductRepository(dbConnection)
	UserRepository := repository.NewUserRepository(dbConnection)
	//Camada usecase
	ProductUsecase := usecase.NewProductUsecase(ProductRepository)
	UserUsecase := usecase.NewUserUsecase(UserRepository)
	//Camada de controller

	ProductController := controller.NewProductController(ProductUsecase)
	UserController := controller.NewUserController(UserUsecase)
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	server.GET("/products", model.VerifyTokenMiddleware, ProductController.GetProducts)
	server.POST("/product", model.VerifyTokenMiddleware, ProductController.CreateProduct)
	server.GET("/product/:productId", model.VerifyTokenMiddleware, ProductController.GetProductById)
	server.PUT("/product/:productId", model.VerifyTokenMiddleware, ProductController.UpdateProductById)
	server.DELETE("/product/:productId", model.VerifyTokenMiddleware, ProductController.DeleteProduct)

	server.POST("/create/user", UserController.CreateUser)
	server.DELETE("/delete/user", model.VerifyTokenMiddleware, UserController.DeleteUser)
	server.POST("/signin", UserController.LogIn)
	server.POST("/refresh_token", UserController.RefreshToken)

	server.Run(":8000")
}
