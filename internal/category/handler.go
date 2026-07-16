package category

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/farrasmumtaz/RentVibe/internal/dto"
	"github.com/farrasmumtaz/RentVibe/internal/models"
	"github.com/farrasmumtaz/RentVibe/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(ctx *gin.Context) {

	var req CreateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.service.Create(&category)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		ctx,
		http.StatusCreated,
		"Category created successfully",
		category,
	)
}

func (h *Handler) FindAll(ctx *gin.Context) {

	search := ctx.DefaultQuery("search", "")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	categories, total, err := h.service.FindAll(search, page, limit)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		responses = append(responses, dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	response.Success(
		ctx,
		http.StatusOK,
		"Categories retrieved successfully",
		gin.H{
			"items": responses,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	)
}

func (h *Handler) FindByID(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid category id")
		return
	}

	category, err := h.service.FindByID(uint(id))
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "Category not found")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result := dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	response.Success(
		ctx,
		http.StatusOK,
		"Category retrieved successfully",
		result,
	)
}
