package item

import (
	"net/http"
	"strconv"

	"github.com/farrasmumtaz/RentVibe/internal/dto"
	"github.com/farrasmumtaz/RentVibe/internal/models"
	"github.com/farrasmumtaz/RentVibe/pkg/response"

	"github.com/gin-gonic/gin"
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

	var req CreateItemRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	item := models.Item{
		Name:        req.Name,
		Description: req.Description,
		PricePerDay: req.PricePerDay,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		CategoryID:  req.CategoryID,
	}

	if err := h.service.Create(&item); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	res := dto.ItemResponse{
		ID:          item.ID,
		CategoryID:  item.CategoryID,
		Name:        item.Name,
		Description: item.Description,
		PricePerDay: item.PricePerDay,
		Stock:       item.Stock,
		ImageURL:    item.ImageURL,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	response.Success(
		ctx,
		http.StatusCreated,
		"Item created successfully",
		res,
	)
}

func (h *Handler) FindAll(ctx *gin.Context) {
	search := ctx.DefaultQuery("search", "")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	items, total, err := h.service.FindAll(search, page, limit)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var result []dto.ItemResponse

	for _, item := range items {
		result = append(result, dto.ItemResponse{
			ID:          item.ID,
			CategoryID:  item.CategoryID,
			Name:        item.Name,
			Description: item.Description,
			PricePerDay: item.PricePerDay,
			Stock:       item.Stock,
			ImageURL:    item.ImageURL,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	response.Success(ctx, http.StatusOK, "Items retrieved successfully", gin.H{
		"items": result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
