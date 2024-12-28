package controller

import (
	"fmt"
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productUsecase usecase.ProductUsecase
}

const (
	ErrorIdCannotBeNull   = "Product id cannot be null"
	ErrorIdNeedToBeNumber = "Product id need to be a number"
	ErrorProductNotFound  = "Product not found"
)

func NewProductController(usecase usecase.ProductUsecase) ProductController {
	return ProductController{
		productUsecase: usecase,
	}
}

func (p *ProductController) GetProducts(ctx *gin.Context) {
	products, err := p.productUsecase.GetProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, products)
}

func (p *ProductController) CreateProduct(ctx *gin.Context) {
	var product model.Product
	err := ctx.BindJSON(&product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	insertedProduct, err := p.productUsecase.CreateProducts(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, insertedProduct)
}

func (p *ProductController) GetProductById(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: ErrorIdCannotBeNull,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: ErrorIdNeedToBeNumber,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUsecase.GetProductById(productId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		response := model.Response{
			Message: ErrorProductNotFound,
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, product)

}

func (p *ProductController) UpdateProductById(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: ErrorIdCannotBeNull,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: ErrorIdNeedToBeNumber,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var product model.Product
	errorBindJson := ctx.BindJSON(&product)
	if errorBindJson != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Printf("id = %d", productId)
	updatedProduct, err := p.productUsecase.UpdateProductById(product, productId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if updatedProduct == nil {
		response := model.Response{
			Message: ErrorProductNotFound,
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	ctx.JSON(http.StatusOK, updatedProduct)

}

func (p *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: ErrorIdCannotBeNull,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: ErrorIdNeedToBeNumber,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUsecase.DeleteProduct(productId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		response := model.Response{
			Message: ErrorProductNotFound,
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	response := model.Response{
		Message: "Product successfully deleted",
	}
	ctx.JSON(http.StatusOK, response)

}
