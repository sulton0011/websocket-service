package main

import (
	"websocket-service/api"
	handler "websocket-service/api/handlers"
	"websocket-service/config"
	"websocket-service/grpc/client"
	"websocket-service/pkg/logger"

	"websocket-service/socket"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	loggerLevel := logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}
	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	c, err := client.NewGrpcClients(cfg)
	if err != nil {
		log.Panic("client.NewGrpcClients", logger.Error(err))
	}

	hub := socket.NewHub(log)
	h := handler.NewHandler(cfg, log, c, hub)
	r := api.SetUpRouter(h, cfg)

	go hub.Run()
	go hub.Read()

	log.Info("HTTP: Server being started...", logger.String("port", cfg.HTTPPort))

	r.Run(cfg.HTTPPort)
}
