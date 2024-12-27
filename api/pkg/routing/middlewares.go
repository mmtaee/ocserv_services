package routing

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
	"sync"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			c.SetRequest(c.Request().WithContext(ctx))

			done := make(chan error, 1)
			go func() {
				done <- next(c)
			}()

			select {
			case err := <-done:
				return err
			case <-ctx.Done():
				return echo.NewHTTPError(http.StatusRequestTimeout, "Request Timeout")
			}
		}
	}
}

var rateLimiters = struct {
	sync.Mutex
	clients map[string]*rate.Limiter
}{
	clients: make(map[string]*rate.Limiter),
}

func getLimiter(k string, r rate.Limit, burst int) *rate.Limiter {
	rateLimiters.Lock()
	defer rateLimiters.Unlock()
	if limited, exists := rateLimiters.clients[k]; exists {
		return limited
	}
	limiter := rate.NewLimiter(r, burst)
	rateLimiters.clients[k] = limiter
	return limiter
}

func calculateRateLimit(count int, per string) (rate.Limit, error) {
	switch strings.ToLower(per) {
	case "s", "seconds", "second":
		return rate.Limit(float64(count)), nil // Per second
	case "m", "minutes", "minute":
		return rate.Limit(float64(count) / 60), nil // Per minute
	case "h", "hours", "hour":
		return rate.Limit(float64(count) / 3600), nil // Per hour
	default:
		return 0, fmt.Errorf("invalid time unit: %s", per)
	}
}

func RateLimitMiddleware(count int, per string, burst int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := fmt.Sprintf("%s:%s", c.Path(), c.RealIP())
			r, err := calculateRateLimit(count, per)
			if err != nil {
				panic(err)
			}
			limiter := getLimiter(key, r, burst)
			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"message": "Rate limit exceeded",
				})
			}
			return next(c)
		}
	}
}
