package category

import "github.com/farrasmumtaz/RentVibe/internal/models"

type Service interface {
	Create(category *models.Category) error
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
