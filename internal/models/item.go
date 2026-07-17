package models

import (
	"time"

	"gorm.io/gorm"
)

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
