package persistence

import (
	"context"
	"errors"
	"fmt"
	"onevideogo/domain"
	"onevideogo/persistence/common"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type IProductRepository interface {
	GetAllProducts() []domain.Product
	GetAllProductsByStore(storeName string) []domain.Product
	AddProduct(product domain.Product) error
	GetById(productId int64) (domain.Product, error)
	DeleteById(productId int64) error
	UpdatePrice(productId int64, newPrice float32) error
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) IProductRepository {
	return &ProductRepository{
		dbPool: dbPool,
	}
}

func (productRepository *ProductRepository) GetAllProducts() []domain.Product {
	ctx := context.Background()
	productRows, err := productRepository.dbPool.Query(ctx, "SELECT * FROM products")

	if err != nil {
		log.Error("Unable to get products: %v\n", err)
		return []domain.Product{}
	}
	return extractProductsFromRows(productRows)
}

func (productRepository *ProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	ctx := context.Background()
	getProductsByStoreNameSql := `SELECT * FROM products WHERE store = $1`
	productRows, err := productRepository.dbPool.Query(ctx, getProductsByStoreNameSql, storeName)

	if err != nil {
		log.Error("Unable to get products: %v\n", err)
		return []domain.Product{}
	}

	return extractProductsFromRows(productRows)
}

func (productRepository *ProductRepository) AddProduct(product domain.Product) error {
	ctx := context.Background()
	addProductSql := `INSERT INTO products (name, price, discount, store) VALUES ($1, $2, $3, $4)`
	addNewProduct, err := productRepository.dbPool.Exec(ctx, addProductSql, product.Name, product.Price, product.Discount, product.Store)

	if err != nil {
		log.Errorf("Unable to add product: %v\n", err)
		return err
	}
	log.Info(fmt.Printf("Product added successfully: %v\n", addNewProduct))
	return nil
}

func (productRepository *ProductRepository) GetById(productId int64) (domain.Product, error) {
	ctx := context.Background()
	getProductByIdSql := `SELECT * FROM products WHERE id = $1`
	productRow := productRepository.dbPool.QueryRow(ctx, getProductByIdSql, productId)

	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	err := productRow.Scan(&id, &name, &price, &discount, &store)
	if err != nil && err.Error() == common.NOT_FOUND {
		log.Errorf("Unable to get product: %v", err)
		return domain.Product{}, errors.New(fmt.Sprintf("Product not found with id %v", productId))
	}

	if err != nil {
		return domain.Product{}, errors.New(fmt.Sprintf("Error while getting id %v", productId))
	}

	return domain.Product{
		Id:       id,
		Name:     name,
		Price:    price,
		Discount: discount,
		Store:    store,
	}, nil
}

func (productRepository *ProductRepository) DeleteById(productId int64) error {
	ctx := context.Background()

	_, getErr := productRepository.GetById(productId)
	if getErr != nil {
		return errors.New(fmt.Sprintf("Product not found with id %v", productId))
	}

	deleteProductByIdSql := `DELETE FROM products WHERE id = $1`
	_, deleteErr := productRepository.dbPool.Exec(ctx, deleteProductByIdSql, productId)
	if deleteErr != nil {
		log.Errorf("Unable to delete product: %v\n", deleteErr)
		return errors.New(fmt.Sprintf("Error while deleting id %v", productId))
	}

	log.Info(fmt.Printf("Product deleted successfully: %v\n", productId))
	return nil
}

func (productRepository *ProductRepository) UpdatePrice(productId int64, newPrice float32) error {
	ctx := context.Background()

	_, getErr := productRepository.GetById(productId)
	if getErr != nil {
		return errors.New(fmt.Sprintf("Product not found with id %v", productId))
	}

	updatePriceSql := `UPDATE products SET price = $1 WHERE id = $2`
	_, err := productRepository.dbPool.Exec(ctx, updatePriceSql, newPrice, productId)

	if err == nil {
		return errors.New(fmt.Sprintf("Error while updating id %v", productId))
	}

	log.Info(fmt.Printf("Product updated successfully: %d, with new price %v", productId, newPrice))
	return nil
}

func extractProductsFromRows(productRows pgx.Rows) []domain.Product {
	var products = []domain.Product{}
	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	for productRows.Next() {
		productRows.Scan(&id, &name, &price, &discount, &store)
		products = append(products, domain.Product{
			Id:       id,
			Name:     name,
			Price:    price,
			Discount: discount,
			Store:    store,
		})
	}
	return products
}
