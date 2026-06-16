package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/XaiPhyr/rdev-go-auth/internal/config"
	"github.com/XaiPhyr/rdev-go-auth/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {
	arg := flag.String("env", "local", "Config environment [local]")
	flag.Parse()

	file := ""
	switch *arg {
	case "local":
		file = "config.yaml"
	}

	cfg, err := config.LoadConfig(file)
	if err != nil {
		log.Println(fmt.Errorf("failed to load config: %w", err))
		return
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	db := config.ConnectDB(cfg.Database)
	server.Container(router, db, cfg)

	router.Run(cfg.Server.Port)
}
