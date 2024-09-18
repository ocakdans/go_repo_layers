package service

import (
	"errors"
	"onevideogo/domain"
	"onevideogo/persistence"
	"onevideogo/service/models"
)

type IProductService interface {
	AllProducts() []domain.Product
	AllProductsByStore(storeName string) []domain.Product
	AddProduct(productCreate models.ProductCreate) error //different from domain.Product as it has no ID, we dont want to care of ID in service layer
	GetById(productId int64) (domain.Product, error)
	DeleteById(productId int64) error
	UpdatePrice(productId int64, newPrice float32) error
}

type ProductService struct {
	productRepository persistence.IProductRepository
}

func NewProductService(productRepository persistence.IProductRepository) IProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (productService *ProductService) AllProducts() []domain.Product {
	return productService.productRepository.GetAllProducts()
}

func (productService *ProductService) AllProductsByStore(storeName string) []domain.Product {
	return productService.productRepository.GetAllProductsByStore(storeName)
}

func (productService *ProductService) AddProduct(productCreate models.ProductCreate) error {
	validateErr := validateProductCreate(productCreate)
	if validateErr != nil {
		return validateErr
	}
	return productService.productRepository.AddProduct(domain.Product{
		Name:     productCreate.Name,
		Price:    productCreate.Price,
		Discount: productCreate.Discount,
		Store:    productCreate.Store,
	})
}

func (productService *ProductService) GetById(productId int64) (domain.Product, error) {
	return productService.productRepository.GetById(productId)
}

func (productService *ProductService) DeleteById(productId int64) error {
	return productService.productRepository.DeleteById(productId)
}

func (productService *ProductService) UpdatePrice(productId int64, newPrice float32) error {
	return productService.productRepository.UpdatePrice(productId, newPrice)
}

func validateProductCreate(productCreate models.ProductCreate) error {
	if productCreate.Discount > 70.0 {
		return errors.New("Discount cannot be more than 70")
	}
	return nil
}
