package category

import (
	"github.com/farrasmumtaz/RentVibe/config"
	"github.com/farrasmumtaz/RentVibe/internal/models"
)

type Repository interface {
	Create(category *models.Category) error
	FindAll(search string, page int, limit int) ([]models.Category, int64, error)
	FindByID(id uint) (*models.Category, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(category *models.Category) error {
	return config.DB.Create(category).Error
}

func (r *repository) FindAll(search string, page int, limit int) ([]models.Category, int64, error) {

	var categories []models.Category
	var total int64

	query := config.DB.Model(&models.Category{})

	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = query.
		Limit(limit).
		Offset(offset).
		Find(&categories).Error

	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *repository) FindByID(id uint) (*models.Category, error) {
	var category models.Category

	err := config.DB.
		Preload("Items").
		First(&category, id).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}
