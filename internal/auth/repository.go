package auth

import (
	"github.com/farrasmumtaz/RentVibe/config"
)

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(user *User) error {
	return config.DB.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
