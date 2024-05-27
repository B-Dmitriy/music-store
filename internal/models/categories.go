package models

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateCategoryBody struct {
	Name string `json:"name" validate:"required,min=3"`
}

type UpdateCategoryBody struct {
	Name string `json:"name" validate:"required,min=3"`
}
