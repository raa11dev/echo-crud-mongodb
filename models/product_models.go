package models

type Product struct {
	Id           int    `json:"id,omitempty" validate:"required"`
	Product_name string `json:"product_name" validate:"required"`
	Status       bool   `json:"status"`
}

type ProductCreate struct {
	Id           int
	Product_name string
	Status       string
}
