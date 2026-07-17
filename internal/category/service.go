package category

import "github.com/farrasmumtaz/RentVibe/internal/models"

type Service interface {
	Create(category *models.Category) error
	FindAll(search string, page int, limit int) ([]models.Category, int64, error)
	FindByID(id uint) (*models.Category, error)
	Update(id uint, req UpdateCategoryRequest) (*models.Category, error)
	Patch(id uint, req PatchCategoryRequest) (*models.Category, error)
	Delete(id uint) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(category *models.Category) error {
	return s.repository.Create(category)
}

func (s *service) FindAll(search string, page int, limit int) ([]models.Category, int64, error) {
	return s.repository.FindAll(search, page, limit)
}

func (s *service) FindByID(id uint) (*models.Category, error) {
	return s.repository.FindByID(id)
}

func (s *service) Update(id uint, req UpdateCategoryRequest) (*models.Category, error) {

	category, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description

	err = s.repository.Update(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *service) Patch(id uint, req PatchCategoryRequest) (*models.Category, error) {

	category, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}

	if req.Description != nil {
		category.Description = *req.Description
	}

	err = s.repository.Patch(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *service) Delete(id uint) error {

	_, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(id)
}
