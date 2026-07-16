package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID         uint     `gorm:"primaryKey"`
	CategoryID uint     `gorm:"not null"`
	Category   Category `gorm:"foreignKey:CategoryID"`

	Name        string  `gorm:"size:150;not null"`
	Description string  `gorm:"type:text"`
	PricePerDay float64 `gorm:"not null"`
	Stock       int     `gorm:"default:0"`
	ImageURL    string  `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
