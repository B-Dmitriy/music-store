package models

type Product struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

type GetAllProductsParams struct {
	Page  int `json:"page" validate:"min=0"`
	Limit int `json:"limit" validate:"min=0"`
}

type ProductsFilter struct {
	CategoryID int `json:"categoryID" validate:"min=0"`
}

type ProductCreateBody struct {
	Title       string  `json:"title" validate:"required,min=3"`
	Price       float64 `json:"price" validate:"required,min=0"`
	CategoryIDs []int   `json:"categoryIDs" validate:"required,min=1"`
}

type ProductUpdateBody struct {
	Title       string  `json:"title" validate:"required,min=3"`
	Price       float64 `json:"price" validate:"required,min=0"`
	CategoryIDs []int   `json:"categoryIDs" validate:"required,min=1"`
}
