package routes

import (
	"github.com/farrasmumtaz/RentVibe/internal/category"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	categoryRepository := category.NewRepository()
	categoryService := category.NewService(categoryRepository)
	categoryHandler := category.NewHandler(categoryService)

	api := router.Group("/api/v1")
	{
		api.POST("/categories", categoryHandler.Create)
	}

	return router
}
