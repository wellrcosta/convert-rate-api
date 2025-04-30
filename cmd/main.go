package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"convert-rate-api/internal/cache"
	"convert-rate-api/internal/handler"
	"convert-rate-api/internal/logger"
	"convert-rate-api/internal/metrics"
	"convert-rate-api/internal/rate"
)

func main() {
	_ = godotenv.Load()

	err := logger.InitLokiLogger(os.Getenv("LOKI_URL"))
	if err != nil {
		return
	}
	defer logger.ShutdownLoki()

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.LogError("Failed to connect to Redis", err)
	}

	cacheService := cache.NewRedisCache(rdb)

	router := gin.Default()

	router.Use(rate.LimiterMiddleware())

	router.GET("/metrics", metrics.PrometheusHandler())
	router.GET("/convert", handler.NewConverterHandler(cacheService))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	logger.LogInfo(fmt.Sprintf("Starting server on port %s", port))
	if err := router.Run(":" + port); err != nil {
		logger.LogError("Failed to start HTTP server", err)
	}
}
