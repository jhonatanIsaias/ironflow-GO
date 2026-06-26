package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)


func RateLimitAuth() gin.HandlerFunc {

	rate, err := limiter.NewRateFromFormatted("5-M")
	if err != nil {
		log.Fatalf("Falha ao configurar a regra de rate limit: %v", err)
	}

	store := memory.NewStore()

	instance := limiter.New(store, rate)

	return mgin.NewMiddleware(instance)
}