package dto

import "time"

type ItemResponse struct {
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

type ItemDetailResponse struct {
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
