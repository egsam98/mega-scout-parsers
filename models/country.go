package models

type Country struct {
	Id   int     `json:"id" example:"141"`
	Name string  `json:"name" example:"Russia"`
	Code *string `json:"code" example:"RU"` // nullable
}
