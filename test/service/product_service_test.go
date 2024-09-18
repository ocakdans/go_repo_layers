package service

import (
	"onevideogo/domain"
	"onevideogo/service"
	"onevideogo/service/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var productService service.IProductService

func TestMain(m *testing.M) {
	initialProducts := []domain.Product{
		{
			Id:       1,
			Name:     "Product 1",
			Price:    100.0,
			Discount: 0.0,
			Store:    "Store 1",
		},
		{
			Id:       2,
			Name:     "Product 2",
			Price:    200.0,
			Discount: 0.0,
			Store:    "Store 2",
		},
	}
	fakeProductRepository := NewFakeProductRepository(initialProducts)
	productService = service.NewProductService(fakeProductRepository)
	os.Exit(m.Run())
}

func Test_ShouldGetAllProducts(t *testing.T) {
	t.Run("GetAllProducts function", func(t *testing.T) {
		actualProducts := productService.AllProducts()
		assert.Equal(t, 2, len(actualProducts))
	})
}

func Test_WhenNoValidationErrorOccurred_ShouldAddProduct(t *testing.T) {
	t.Run("WhenNoValidationErrorOccurred_ShouldAddProduct function", func(t *testing.T) {
		productService.AddProduct(models.ProductCreate{
			Name:     "Product 3",
			Price:    300.0,
			Discount: 0.0,
			Store:    "Store 3",
		})
		actualProducts := productService.AllProducts()
		assert.Equal(t, 3, len(actualProducts))
	})

}

func Test_WhenValidationErrorOccurred_ShouldNotAddProduct(t *testing.T) {
	t.Run("WhenNoValidationErrorOccurred_ShouldAddProduct function", func(t *testing.T) {
		err := productService.AddProduct(models.ProductCreate{
			Name:     "Product 4",
			Price:    300.0,
			Discount: 80.0,
			Store:    "Store 4",
		})
		actualProducts := productService.AllProducts()
		assert.Equal(t, 2, len(actualProducts))
		assert.Equal(t, "Discount cannot be more than 70", err.Error())
	})

}
