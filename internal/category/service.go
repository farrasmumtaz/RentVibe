package category

import "github.com/farrasmumtaz/RentVibe/internal/models"

type Service interface {
	Create(category *models.Category) error
	FindAll(search string, page int, limit int) ([]models.Category, int64, error)
	FindByID(id uint) (*models.Category, error)
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
