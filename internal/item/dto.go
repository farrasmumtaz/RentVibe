package item

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
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type PatchItemRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
	CategoryID  *uint    `json:"category_id,omitempty"`
}
