package main

import (
	"fmt"
	"log"

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

	db := config.ConnectDB(cfg.Database)
	redis := config.ConnectRedis(cfg.Redis)
	server.Container(router, db, redis, cfg)

	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
