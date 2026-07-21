package auth

import (
	"github.com/farrasmumtaz/RentVibe/config"
	"github.com/farrasmumtaz/RentVibe/internal/models"
)

type Repository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
