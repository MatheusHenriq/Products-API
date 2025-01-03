package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type ProductUsecase struct {
	repository repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return ProductUsecase{repository: repo}
}

func (pu *ProductUsecase) GetProducts() ([]model.Product, error) {
	return pu.repository.GetProducts()
}

func (pu *ProductUsecase) CreateProducts(product model.Product) (model.Product, error) {
	productId, err := pu.repository.CreateProduct(product)
	if err != nil {
		return model.Product{}, err
	}
	product.ID = productId
	return product, nil
}

func (pu *ProductUsecase) GetProductById(id_product int) (*model.Product, error) {
	product, err := pu.repository.GetProductById(id_product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pu *ProductUsecase) UpdateProductById(p model.Product, pId int) (*model.Product, error) {
	product, err := pu.repository.UpdateProductById(p, pId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pu *ProductUsecase) DeleteProduct(pId int) (*model.Product, error) {
	product, err := pu.repository.DeleteProduct(pId)
	if err != nil {
		return nil, err
	}
	return product, nil
}
