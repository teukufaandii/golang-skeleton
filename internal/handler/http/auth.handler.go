package handler

import (
	"golang-skeleton/internal/dto/request"
	"golang-skeleton/internal/service"
	"golang-skeleton/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.Validate(req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	result, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if err == services.ErrEmailExists {
			utils.Conflict(c, "Email already registered")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result)
}
