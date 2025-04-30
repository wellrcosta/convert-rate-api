package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"exchange-rate-api/internal/cache"
	"exchange-rate-api/internal/handler"
	"exchange-rate-api/internal/metrics"
	"exchange-rate-api/internal/rate"
)

func main() {
	// Load .env file if exists
	_ = godotenv.Load()

	// Set up logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Read Redis configuration
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	// Initialize cache layer
	cacheService := cache.NewRedisCache(rdb)

	// Set up Gin and routes
	router := gin.Default()

	// Apply rate limiting middleware
	router.Use(rate.LimiterMiddleware())

	// Register metrics endpoint
	router.GET("/metrics", metrics.PrometheusHandler())

	// Register currency conversion handler
	router.GET("/convert", handler.NewConverterHandler(cacheService))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Info().Msgf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
