package item

import "time"

type CreateItemRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Description string  `json:"description"`
	PricePerDay float64 `json:"price_per_day" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	ImageURL    string  `json:"image_url"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type UpdateItemRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Description string  `json:"description"`
	PricePerDay float64 `json:"price_per_day" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	ImageURL    string  `json:"image_url"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type PatchItemRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description *string  `json:"description,omitempty"`
	PricePerDay *float64 `json:"price_per_day,omitempty" binding:"omitempty,gt=0"`
	Stock       *int     `json:"stock,omitempty" binding:"omitempty,gte=0"`
	ImageURL    *string  `json:"image_url,omitempty"`
	CategoryID  *uint    `json:"category_id,omitempty" binding:"omitempty,gt=0"`
}

type Response struct {
	ID          uint      `json:"id"`
	CategoryID  uint      `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PricePerDay float64   `json:"price_per_day"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DetailResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	PricePerDay float64          `json:"price_per_day"`
	Stock       int              `json:"stock"`
	ImageURL    string           `json:"image_url"`
	Category    CategoryResponse `json:"category"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
