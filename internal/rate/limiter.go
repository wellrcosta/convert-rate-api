package rate

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginMiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memoryStore "github.com/ulule/limiter/v3/drivers/store/memory"
	"os"
)

// LimiterMiddleware returns a Gin middleware that applies rate limiting based on IP.
func LimiterMiddleware() gin.HandlerFunc {
	limit := os.Getenv("RATE_LIMIT_REQUESTS")
	window := os.Getenv("RATE_LIMIT_WINDOW")

	if limit == "" {
		limit = "100"
	}
	if window == "" {
		window = "15m"
	}

	// Format: "100-15m" means 100 requests per 15 minutes
	rateStr := limit + "-" + window

	rateConfig, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		panic(err)
	}

	store := memoryStore.NewStore()
	instance := limiter.New(store, rateConfig)

	return ginMiddleware.NewMiddleware(instance)
}
