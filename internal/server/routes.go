package server

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api-template/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api-template/internal/auth"
	"github.com/XaiPhyr/rdev-go-api-template/internal/config"
	"github.com/XaiPhyr/rdev-go-api-template/internal/middleware"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api-template/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

func Container(r *gin.Engine, db *bun.DB, redis *redis.Client, cfg *config.Config) {
	emailSvc := email.NewEmailService(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.From)

	auditLogRepo := audit_logs.NewAuditLogRepository(db)
	auditLogSvc := audit_logs.NewAuditLogService(auditLogRepo)
	go auditLogSvc.QueAuditLog(context.Background())

	authRepo := auth.NewAuthRepository(db)
	authSvc := auth.NewAuthService(authRepo, cfg, emailSvc, redis)
	userRepo := users.NewUserRepository(db)

	apiVersion := r.Group(cfg.Server.Version)
	apiVersion.Use(middleware.RateLimiter())

	setupAuthRoutes(apiVersion, authSvc)
	setupUserRoutes(apiVersion, userRepo, authSvc, emailSvc, redis)
}

func setupAuthRoutes(rg *gin.RouterGroup, authSvc auth.AuthService) {
	h := auth.NewAuthHandler(authSvc)

	rg.POST("/login", h.Login)
	rg.POST("/register", h.Register)
}

func setupUserRoutes(rg *gin.RouterGroup, repo users.UserRepository, authSvc auth.AuthService, es email.EmailService, redis *redis.Client) {
	svc := users.NewUserService(repo, es, redis)
	h := users.NewUserHandler(svc)

	route := rg.Group("/service_types")
	route.Use(middleware.AuthRequired(authSvc))

	route.GET("/:uuid", middleware.PermissionRequired(authSvc, "users:read"), h.ReadOne)
	route.GET("", middleware.PermissionRequired(authSvc, "users:read"), h.ReadAll)
	route.POST("", middleware.PermissionRequired(authSvc, "users:create"), h.Create)
	route.PUT("/:uuid", middleware.PermissionRequired(authSvc, "users:update"), h.Update)
	route.DELETE("/:uuid", middleware.PermissionRequired(authSvc, "users:delete"), h.Delete)
}
