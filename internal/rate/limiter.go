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

	if limit == "" {
		limit = "100-M"
	}

	rateConfig, err := limiter.NewRateFromFormatted(limit)
	if err != nil {
		panic(err)
	}

	store := memoryStore.NewStore()
	instance := limiter.New(store, rateConfig)

	return ginMiddleware.NewMiddleware(instance)
}
