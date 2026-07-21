package auth

import (
	"errors"
	"net/http"

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

func (h *Handler) Register(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyUsed) {
			response.Error(ctx, http.StatusConflict, "Email sudah digunakan")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "Registrasi berhasil", result)
}

func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredential) {
			response.Error(ctx, http.StatusUnauthorized, "Email atau password tidak valid")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Login berhasil", result)
}
