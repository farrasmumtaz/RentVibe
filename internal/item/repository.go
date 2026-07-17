package item

import (
	"github.com/farrasmumtaz/RentVibe/config"
	"github.com/farrasmumtaz/RentVibe/internal/models"
)

type Repository interface {
	Create(item *models.Item) error
	FindAll(search string, page int, limit int) ([]models.Item, int64, error)
	FindByID(id uint) (*models.Item, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(item *models.Item) error {
	return config.DB.Create(item).Error
}

func (r *repository) FindAll(search string, page int, limit int) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	query := config.DB.Model(&models.Item{})

	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) FindByID(id uint) (*models.Item, error) {
	var item models.Item

	err := config.DB.
		Debug().
		Preload("Category").
		First(&item, id).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}
