package main

import (
	"fmt"
	"log"

	"github.com/Delkira544/rakiduam/config"
	"github.com/Delkira544/rakiduam/internal/app"
	"github.com/Delkira544/rakiduam/internal/platform/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := logger.New(cfg.Log); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.L().Sync()

	ginMode := gin.ReleaseMode
	if cfg.App.Env == "development" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	logger.L().Info("starting server",
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.Int("port", cfg.App.Port))

	a, cleanup, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer cleanup()

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("%s starting on %s", cfg.App.Name, addr)
	if err := a.Router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
