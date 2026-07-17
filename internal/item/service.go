package item

import "github.com/farrasmumtaz/RentVibe/internal/models"

type Service interface {
	Create(item *models.Item) error
	FindAll(search string, page int, limit int) ([]models.Item, int64, error)
	FindByID(id uint) (*models.Item, error)
	Update(id uint, item *models.Item) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(item *models.Item) error {
	return s.repository.Create(item)
}

func (s *service) FindAll(search string, page int, limit int) ([]models.Item, int64, error) {
	return s.repository.FindAll(search, page, limit)
}

func (s *service) FindByID(id uint) (*models.Item, error) {
	return s.repository.FindByID(id)
}

func (s *service) Update(id uint, item *models.Item) error {
	return s.repository.Update(id, item)
}
