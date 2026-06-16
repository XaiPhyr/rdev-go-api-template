package users

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/helpers"
	"github.com/gin-gonic/gin"
)

type handler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *handler {
	return &handler{svc: svc}
}

func (h *handler) Create(ctx *gin.Context) {
	var req UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": helpers.ParseValidationErr(err)})
		return
	}

	err := h.svc.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (h *handler) ReadOne(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	if uuid == "" || len(uuid) < 36 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	user, err := h.svc.ReadOne(ctx.Request.Context(), uuid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": user})
}

func (h *handler) ReadAll(ctx *gin.Context) {
	var req dto.BaseFilters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": helpers.ParseValidationErr(err)})
		return
	}

	users, total, err := h.svc.ReadAll(ctx.Request.Context(), req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "total": total, "data": users})
}

func (h *handler) Update(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" || len(uuid) < 36 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	var req UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": helpers.ParseValidationErr(err)})
		return
	}

	if err := h.svc.Update(ctx.Request.Context(), uuid, req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *handler) Delete(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	if uuid == "" || len(uuid) < 36 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	if err := h.svc.Delete(ctx.Request.Context(), uuid); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
