package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XaiPhyr/rdev-go-api-template/internal/config"
	"github.com/XaiPhyr/rdev-go-api-template/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Println(fmt.Errorf("failed to load config: %w", err))
		return
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.MaxMultipartMemory = 8 << 20

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	db := config.ConnectDB(cfg.Database)
	redis := config.ConnectRedis(cfg.Redis)
	server.Container(router, db, redis, cfg)

	srv := &http.Server{
		Addr:              cfg.Server.Port,
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server initiating listening routines on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Listen error encountered: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("⚠️  Termination signal received (%s). Triggering graceful teardown...\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Forced server shutdown triggered: %v", err)
	}

	log.Println("🗄️  Closing Bun database connection pools...")
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("⚠️  Error closing database: %v\n", err)
		}
	}

	log.Println("⚡ Closing Redis client connections...")
	if redis != nil {
		if err := redis.Close(); err != nil {
			log.Printf("⚠️  Error closing Redis cache client: %v\n", err)
		}
	}

	log.Println("✨ Server instance successfully wrapped down and exited cleanly.")
}
