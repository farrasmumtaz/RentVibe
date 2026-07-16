package routes

import (
	"github.com/farrasmumtaz/RentVibe/internal/category"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	api := router.Group("/api/v1")

	// Category
	{
		categoryRepository := category.NewRepository()
		categoryService := category.NewService(categoryRepository)
		categoryHandler := category.NewHandler(categoryService)

		api.POST("/categories", categoryHandler.Create)
		api.GET("/categories", categoryHandler.FindAll)
		api.GET("/categories/:id", categoryHandler.FindByID)
		api.PUT("/categories/:id", categoryHandler.Update)
		api.PATCH("/categories/:id", categoryHandler.Patch)
	}
	return router
}
