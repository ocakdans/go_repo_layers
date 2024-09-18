package controller

import (
	"net/http"
	"onevideogo/controller/request"
	"onevideogo/controller/response"
	"onevideogo/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productService service.IProductService
}

func NewProductController(productService service.IProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (productController *ProductController) RegisterRoutes(e *echo.Echo) {

	e.GET("api/v1/products/:id", productController.GetProductById)
	e.GET("api/v1/products", productController.GetAllProducts)
	e.POST("api/v1/products", productController.AddProduct)
	e.PUT("api/v1/products/:id", productController.UpdatePrice)
	e.DELETE("api/v1/products/:id", productController.DeleteProductById)

}

func (productController *ProductController) GetProductById(c echo.Context) error {
	productIdStr := c.Param("id")
	productId, _ := strconv.Atoi(productIdStr)
	product, err := productController.productService.GetById(int64(productId))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response.ToResponse(product))

}

func (productController *ProductController) GetAllProducts(c echo.Context) error {
	store := c.QueryParam("store")
	if len(store) == 0 {
		allProducts := productController.productService.AllProducts()
		return c.JSON(http.StatusOK, response.ToResponseList(allProducts))
	} else {
		productsWithStore := productController.productService.AllProductsByStore(store)
		return c.JSON(http.StatusOK, response.ToResponseList(productsWithStore))
	}
}

func (productController *ProductController) AddProduct(c echo.Context) error {
	var addProductRequest request.AddProductRequest
	err := c.Bind(&addProductRequest) // request body deki değerleri addProductRequest nesnesine bind eder, map eder.
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}
	err = productController.productService.AddProduct(addProductRequest.ToModel())
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})

	}
	return c.NoContent(http.StatusCreated)

}

func (productController *ProductController) UpdatePrice(c echo.Context) error {
	productIdStr := c.Param("id")
	productId, _ := strconv.Atoi(productIdStr)
	newPrice := c.QueryParam("newPrice")
	if len(newPrice) == 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "newPrice query parameter is required",
		})
	}

	convertedSNewPrice, err := strconv.ParseFloat(newPrice, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}
	productController.productService.UpdatePrice(int64(productId), float32(convertedSNewPrice))
	return c.NoContent(http.StatusOK)

}

func (productController *ProductController) DeleteProductById(c echo.Context) error {
	c.Param("id")
	productIdStr := c.Param("id")
	productId, _ := strconv.Atoi(productIdStr)
	err := productController.productService.DeleteById(int64(productId))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})

	}
	return c.NoContent(http.StatusOK)
}
