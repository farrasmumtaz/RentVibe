package migration

import (
	"fmt"

	"github.com/farrasmumtaz/RentVibe/config"
	"github.com/farrasmumtaz/RentVibe/internal/auth"
	"github.com/farrasmumtaz/RentVibe/internal/catalog"
)

func Run() error {
	if err := config.DB.AutoMigrate(
		&auth.User{},
		&catalog.Category{},
		&catalog.Item{},
	); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}
