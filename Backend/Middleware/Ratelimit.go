package Middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"elucid503/SproutLoggingRelay/Types"
)

type RouteRateLimitConfig struct {
	MaxRequestsAllowed int
	TimeWindow         int // in seconds
}

type RateLimitStore struct {
	UUID                 string
	RequestsWithinWindow int
	WindowStart          int64
}

// Initialize the map at package level

var RateLimitStoreMap = make(map[string]RateLimitStore)

func RateLimit(RouteConfig *RouteRateLimitConfig) gin.HandlerFunc {

	return func(Context *gin.Context) {

		UUID := Context.ClientIP()

		ShouldBlock := CheckAndUpdateRateLimit(UUID, RouteConfig.TimeWindow, RouteConfig.MaxRequestsAllowed)

		if ShouldBlock {

			Context.JSON(429, Types.Response{Message: "Rate limit exceeded"})
			Context.Abort()

			return

		}

		Context.Next()

	}

}

func CheckAndUpdateRateLimit(UUID string, TimeWindow int, MaxRequests int) bool {

	Now := time.Now().Unix()

	if Stored, Exists := RateLimitStoreMap[UUID]; Exists {

		if Now-Stored.WindowStart >= int64(TimeWindow) {

			ResetRateLimit(UUID, Now)

			return false

		}

		return IncrementRateLimit(UUID, MaxRequests)

	}

	// First request from this UUID

	InitRateLimitStore(UUID, Now)

	return false
}

func ResetRateLimit(UUID string, now int64) {

	RateLimitStoreMap[UUID] = RateLimitStore{
		UUID:                 UUID,
		RequestsWithinWindow: 1,
		WindowStart:          now,
	}

}

func IncrementRateLimit(UUID string, MaxRequestsAllowed int) bool {

	Stored := RateLimitStoreMap[UUID]
	Stored.RequestsWithinWindow++

	RateLimitStoreMap[UUID] = Stored

	return Stored.RequestsWithinWindow > MaxRequestsAllowed

}

func InitRateLimitStore(UUID string, now int64) {

	RateLimitStoreMap[UUID] = RateLimitStore{
		UUID:                 UUID,
		RequestsWithinWindow: 1,
		WindowStart:          now,
	}

}
