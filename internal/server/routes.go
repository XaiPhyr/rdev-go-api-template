package server

import (
	"github.com/XaiPhyr/rdev-go-auth/internal/auth"
	"github.com/XaiPhyr/rdev-go-auth/internal/config"
	"github.com/XaiPhyr/rdev-go-auth/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func Container(r *gin.Engine, db *bun.DB, cfg *config.Config) {
	userRepo := users.NewUserRepository(db)
	authSvc := auth.NewAuthService(userRepo, cfg)

	apiVersion := r.Group("/api/v1")
	setupAuthRoutes(apiVersion, authSvc)
}

func setupAuthRoutes(rg *gin.RouterGroup, authSvc *auth.AuthService) {
	authHandler := auth.NewAuthHandler(authSvc)

	rg.POST("/login", authHandler.Login)
}
