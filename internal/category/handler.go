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

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

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
		"Kategori berhasil dibuat",
		category,
	)
}

func (h *Handler) FindAll(ctx *gin.Context) {

	search := ctx.DefaultQuery("search", "")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", strconv.Itoa(defaultPageSize)))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
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
		"Kategori berhasil diambil",
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
		response.Error(ctx, http.StatusBadRequest, "ID kategori tidak valid")
		return
	}

	category, err := h.service.FindByID(uint(id))
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "Kategori tidak ditemukan")
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
	for _, item := range category.Items {
		result.Items = append(result.Items, dto.ItemResponse{
			ID: item.ID, CategoryID: item.CategoryID, Name: item.Name,
			Description: item.Description, PricePerDay: item.PricePerDay,
			Stock: item.Stock, ImageURL: item.ImageURL,
			CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		})
	}

	response.Success(
		ctx,
		http.StatusOK,
		"Kategori berhasil diambil",
		result,
	)
}

func (h *Handler) Update(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "ID kategori tidak valid")
		return
	}

	var req UpdateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.Update(uint(id), req)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "Kategori tidak ditemukan")
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
		"Kategori berhasil diperbarui",
		result,
	)
}

func (h *Handler) Patch(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "ID kategori tidak valid")
		return
	}

	var req PatchCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.Patch(uint(id), req)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "Kategori tidak ditemukan")
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
		"Kategori berhasil diperbarui sebagian",
		result,
	)
}

func (h *Handler) Delete(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "ID kategori tidak valid")
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, "Kategori tidak ditemukan")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"Kategori berhasil dihapus",
		nil,
	)
}
