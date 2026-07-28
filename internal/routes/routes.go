package routes

import (
	"net/http"
	"time"

	"github.com/farrasmumtaz/RentVibe/config"
	"github.com/farrasmumtaz/RentVibe/internal/auth"
	"github.com/farrasmumtaz/RentVibe/internal/cache"
	"github.com/farrasmumtaz/RentVibe/internal/category"
	"github.com/farrasmumtaz/RentVibe/internal/item"
	"github.com/farrasmumtaz/RentVibe/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cacheStore cache.Store) (*gin.Engine, error) {

	router := gin.Default()
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(20, time.Minute))

	api := router.Group("/api/v1")

	router.GET("/health", func(ctx *gin.Context) {
		sqlDB, err := config.DB.DB()
		if err != nil || sqlDB.PingContext(ctx.Request.Context()) != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	tokenService, err := auth.NewTokenService()
	if err != nil {
		return nil, err
	}

	authRepository := auth.NewRepository()
	authService := auth.NewService(authRepository, tokenService)
	authHandler := auth.NewHandler(authService)

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	protected := api.Group("")
	protected.Use(middleware.Auth(tokenService))

	// Category
	{
		categoryRepository := category.NewRepository()
		categoryService := category.NewService(categoryRepository, cacheStore)
		categoryHandler := category.NewHandler(categoryService)

		protected.POST("/categories", categoryHandler.Create)
		protected.GET("/categories", categoryHandler.FindAll)
		protected.GET("/categories/:id", categoryHandler.FindByID)
		protected.PUT("/categories/:id", categoryHandler.Update)
		protected.PATCH("/categories/:id", categoryHandler.Patch)
		protected.DELETE("/categories/:id", categoryHandler.Delete)
	}

	{
		itemRepository := item.NewRepository()
		itemService := item.NewService(itemRepository, cacheStore)
		itemHandler := item.NewHandler(itemService)

		protected.POST("/items", itemHandler.Create)
		protected.GET("/items", itemHandler.FindAll)
		protected.GET("/items/:id", itemHandler.FindByID)
		protected.PUT("/items/:id", itemHandler.Update)
		protected.PATCH("/items/:id", itemHandler.Patch)
		protected.DELETE("/items/:id", itemHandler.Delete)

	}
	return router, nil
}
