package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"exchange-rate-api/internal/cache"
	"exchange-rate-api/internal/client"
	"exchange-rate-api/internal/metrics"
)

func NewConverterHandler(c cache.CacheService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		from := ctx.Query("from")
		to := ctx.Query("to")
		amountStr := ctx.Query("amount")

		if from == "" || to == "" || amountStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters: from, to, amount"})
			return
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
			return
		}

		cacheKey := fmt.Sprintf("rate:%s:%s", from, to)
		rateStr, found := c.Get(ctx, cacheKey)

		var rate float64
		var cached bool

		if found {
			rate, _ = strconv.ParseFloat(rateStr, 64)
			cached = true
			metrics.CacheHits.Inc()
		} else {
			apiKey := os.Getenv("EXCHANGE_API_KEY")
			apiURL := os.Getenv("EXCHANGE_API_URL")
			rate, err = client.FetchRate(apiURL, apiKey, from, to)
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch exchange rate")
				metrics.Errors.Inc()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exchange rate"})
				return
			}

			ttlSeconds, _ := strconv.Atoi(os.Getenv("CACHE_TTL_SECONDS"))
			c.Set(ctx, cacheKey, fmt.Sprintf("%.6f", rate), time.Duration(ttlSeconds)*time.Second)
			cached = false
			metrics.CacheMisses.Inc()
		}

		converted := amount * rate
		metrics.Conversions.WithLabelValues(from, to).Inc()

		ctx.JSON(http.StatusOK, gin.H{
			"from":            from,
			"to":              to,
			"amount":          fmt.Sprintf("%.2f", amount),
			"rate":            fmt.Sprintf("%.6f", rate),
			"convertedAmount": fmt.Sprintf("%.2f", converted),
			"cached":          cached,
		})
	}
}
