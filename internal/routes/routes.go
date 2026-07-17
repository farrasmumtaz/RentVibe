package routes

import (
	"github.com/farrasmumtaz/RentVibe/internal/category"
	"github.com/farrasmumtaz/RentVibe/internal/item"

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
		api.DELETE("/categories/:id", categoryHandler.Delete)
	}

	{
		itemRepository := item.NewRepository()
		itemService := item.NewService(itemRepository)
		itemHandler := item.NewHandler(itemService)

		api.POST("/items", itemHandler.Create)
		api.GET("/items", itemHandler.FindAll)
		api.GET("/items/:id", itemHandler.FindByID)
		api.PUT("/items/:id", itemHandler.Update)
		api.PATCH("/items/:id", itemHandler.Patch)

	}
	return router
}
