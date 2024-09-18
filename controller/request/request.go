package request

import "onevideogo/service/models"

type AddProductRequest struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Discount float32 `json:"discount"`
	Store    string  `json:"store"`
}

func (addProductRequest AddProductRequest) ToModel() models.ProductCreate {
	return models.ProductCreate{
		Name:     addProductRequest.Name,
		Price:    addProductRequest.Price,
		Discount: addProductRequest.Discount,
		Store:    addProductRequest.Store,
	}
}
