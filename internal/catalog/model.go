package catalog

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`

	Items []Item `gorm:"foreignKey:CategoryID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Item struct {
	ID uint `gorm:"primaryKey"`

	CategoryID uint
	Category   *Category `gorm:"foreignKey:CategoryID;references:ID"`

	Name        string
	Description string
	PricePerDay float64
	Stock       int
	ImageURL    string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
