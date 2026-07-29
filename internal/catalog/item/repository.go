package item

import (
	"github.com/farrasmumtaz/RentVibe/config"
	"github.com/farrasmumtaz/RentVibe/internal/catalog"
)

type Repository interface {
	Create(item *catalog.Item) error
	FindAll(search string, page int, limit int) ([]catalog.Item, int64, error)
	FindByID(id uint) (*catalog.Item, error)
	Update(id uint, item *catalog.Item) error
	Patch(id uint, req PatchItemRequest) (*catalog.Item, error)
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(item *catalog.Item) error {
	return config.DB.Create(item).Error
}

func (r *repository) FindAll(search string, page int, limit int) ([]catalog.Item, int64, error) {
	var items []catalog.Item
	var total int64

	query := config.DB.Model(&catalog.Item{})

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

func (r *repository) FindByID(id uint) (*catalog.Item, error) {
	var item catalog.Item

	err := config.DB.
		Preload("Category").
		First(&item, id).Error

	return &item, err
}

func (r *repository) Update(id uint, item *catalog.Item) error {
	var existing catalog.Item

	if err := config.DB.First(&existing, id).Error; err != nil {
		return err
	}

	existing.Name = item.Name
	existing.Description = item.Description
	existing.PricePerDay = item.PricePerDay
	existing.Stock = item.Stock
	existing.ImageURL = item.ImageURL
	existing.CategoryID = item.CategoryID

	if err := config.DB.Save(&existing).Error; err != nil {
		return err
	}

	*item = existing
	return nil
}

func (r *repository) Patch(id uint, req PatchItemRequest) (*catalog.Item, error) {
	var item catalog.Item

	if err := config.DB.First(&item, id).Error; err != nil {
		return nil, err
	}

	if req.Name != nil {
		item.Name = *req.Name
	}

	if req.Description != nil {
		item.Description = *req.Description
	}

	if req.PricePerDay != nil {
		item.PricePerDay = *req.PricePerDay
	}

	if req.Stock != nil {
		item.Stock = *req.Stock
	}

	if req.ImageURL != nil {
		item.ImageURL = *req.ImageURL
	}

	if req.CategoryID != nil {
		item.CategoryID = *req.CategoryID
	}

	if err := config.DB.Save(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *repository) Delete(id uint) error {
	return config.DB.Delete(&catalog.Item{}, id).Error
}
