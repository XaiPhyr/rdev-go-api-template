package auth

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/helpers"
	"github.com/gin-gonic/gin"
)

type handler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) *handler {
	return &handler{svc: svc}
}

func (h *handler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	token, err := h.svc.Login(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *handler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": helpers.ParseValidationErr(err)})
		return
	}

	err := h.svc.Register(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
}
