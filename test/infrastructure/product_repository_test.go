package infrastructure

import (
	"context"
	"fmt"
	"onevideogo/common/postgresql"
	"onevideogo/domain"
	"onevideogo/persistence"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

/* func TestAdd(t *testing.T) {
	t.Run("Add function", func(t *testing.T) {
		actual := Add(1, 2)
		expected := 3
		assert.Equal(t, expected, actual)
	})

}

func Add(x, y int) int {
	return x + y
}
*/

var productRepository persistence.IProductRepository
var dbPool *pgxpool.Pool
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()
	dbPool = postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		UserName:              "postgres",
		Password:              "postgres",
		DbName:                "productapp",
		MaxConnections:        "10",
		MaxConnectionIdleTime: "30s",
	})

	productRepository = persistence.NewProductRepository(dbPool)
	fmt.Println("Before running tests")
	exitCode := m.Run()
	fmt.Println("After running tests")
	os.Exit(exitCode)
}

func setup(ctx context.Context, dbPool *pgxpool.Pool) {
	TestDataInitialize(ctx, dbPool)
}

func clear(ctx context.Context, dbPool *pgxpool.Pool) {
	TruncateTestData(ctx, dbPool)
}

func TestGetAllProducts(t *testing.T) {
	setup(ctx, dbPool)

	expectedProducts := []domain.Product{
		{
			Id:       1,
			Name:     "AirFryer",
			Price:    3000.0,
			Discount: 22.0,
			Store:    "ABC TECH",
		},
		{
			Id:       2,
			Name:     "Ütü",
			Price:    1500.0,
			Discount: 10.0,
			Store:    "ABC TECH",
		},
		{
			Id:       3,
			Name:     "Çamaşır Makinesi",
			Price:    10000.0,
			Discount: 15.0,
			Store:    "ABC TECH",
		},
		{
			Id:       4,
			Name:     "Lambader",
			Price:    2000.0,
			Discount: 0.0,
			Store:    "Dekorasyon Sarayı",
		},
	}

	t.Run("GetAllProducts function", func(t *testing.T) {
		actualProducts := productRepository.GetAllProducts()
		time.Sleep(1 * time.Second)
		fmt.Println("Actual Products: ", actualProducts)
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
	})

	clear(ctx, dbPool)

}

func TestGetAllProductsByStore(t *testing.T) {
	setup(ctx, dbPool)

	expectedProducts := []domain.Product{
		{
			Id:       1,
			Name:     "AirFryer",
			Price:    3000.0,
			Discount: 22.0,
			Store:    "ABC TECH",
		},
		{
			Id:       2,
			Name:     "Ütü",
			Price:    1500.0,
			Discount: 10.0,
			Store:    "ABC TECH",
		},
		{
			Id:       3,
			Name:     "Çamaşır Makinesi",
			Price:    10000.0,
			Discount: 15.0,
			Store:    "ABC TECH",
		},
	}
	t.Run("GetAllProductsByStore function", func(t *testing.T) {
		actualProducts := productRepository.GetAllProductsByStore("ABC TECH")
		//fmt.Println("Actual Products: ", actualProducts)
		assert.Equal(t, 3, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
	})
	clear(ctx, dbPool)
}

func TestAddProduct(t *testing.T) {

	expectedProducts := []domain.Product{
		{
			Id:       1,
			Name:     "Kupa",
			Price:    100.0,
			Discount: 0.0,
			Store:    "Kırtasiye Merkezi",
		},
	}
	newProduct := domain.Product{
		Name:     "Kupa",
		Price:    100.0,
		Discount: 0.0,
		Store:    "Kırtasiye Merkezi",
	}
	t.Run("AddProduct function", func(t *testing.T) {
		productRepository.AddProduct(newProduct)
		//fmt.Println("Actual Products: ", actualProducts)
		actualProducts := productRepository.GetAllProducts()
		assert.Equal(t, 1, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
	})
	clear(ctx, dbPool)
}

func TestGetById(t *testing.T) {
	expectedProduct := domain.Product{
		Id:       1,
		Name:     "Kupa",
		Price:    100.0,
		Discount: 0.0,
		Store:    "Kırtasiye Merkezi",
	}
	newProduct := domain.Product{
		Name:     "Kupa",
		Price:    100.0,
		Discount: 0.0,
		Store:    "Kırtasiye Merkezi",
	}
	t.Run("TestGetById function", func(t *testing.T) {
		productRepository.AddProduct(newProduct)
		//fmt.Println("Actual Products: ", actualProducts)
		actualProduct, _ := productRepository.GetById(1)
		_, err := productRepository.GetById(599)
		//assert.Equal(t, 1, len(actualProduct))
		assert.Equal(t, expectedProduct, actualProduct)
		assert.NotNil(t, err)
		assert.Equal(t, "Product not found with id 599", err.Error())
	})
	clear(ctx, dbPool)
}

func TestDeleteById(t *testing.T) {
	setup(ctx, dbPool)
	t.Run("TestDeleteById function", func(t *testing.T) {
		productRepository.DeleteById(1)
		actualProducts := productRepository.GetAllProducts()
		_, err := productRepository.GetById(1)
		assert.Equal(t, 3, len(actualProducts))
		assert.NotNil(t, err)
		assert.Equal(t, "Product not found with id 1", err.Error())
	})
	clear(ctx, dbPool)
}

func TestUpdatePrice(t *testing.T) {
	setup(ctx, dbPool)
	t.Run("TestUpdatePrice function", func(t *testing.T) {
		productBeforeUpdate, _ := productRepository.GetById(1)
		assert.Equal(t, float32(3000.0), productBeforeUpdate.Price)
		productRepository.UpdatePrice(1, 4000.0)
		productAfterUpdate, _ := productRepository.GetById(1)
		assert.Equal(t, float32(4000.0), productAfterUpdate.Price)

	})
	clear(ctx, dbPool)
}
