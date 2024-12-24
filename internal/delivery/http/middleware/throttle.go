package middleware

import (
	"fmt"
	"time"

	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisRateLimiter struct {
	Limiter *redisrate.Limiter
	Rate    int
	Burst   int
	Period  int
}

const RateRequest = "rate_request_%s"

func NewRedisRateLimiter(client *redis.Client, viper *viper.Viper) *RedisRateLimiter {
	rate := viper.GetInt("limiter.rate")
	burst := viper.GetInt("limiter.burst")
	period := viper.GetInt("limiter.period")

	return &RedisRateLimiter{
		Limiter: redisrate.NewLimiter(client),
		Rate:    rate,
		Burst:   burst,
		Period:  period}

}

func (r *RedisRateLimiter) ThrottleByKey(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		res, _ := r.Limiter.Allow(ctx, fmt.Sprintf(RateRequest, key), redisrate.Limit{
			Rate:   r.Rate,
			Burst:  r.Burst,
			Period: time.Duration(r.Period) * time.Minute,
		})
		if res.Allowed <= 0 {
			return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded")
		}
		return c.Next()
	}
}

func (r *RedisRateLimiter) ThrottleByIp() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		ip := c.IP()
		res, _ := r.Limiter.Allow(ctx, fmt.Sprintf(RateRequest, ip), redisrate.Limit{
			Rate:   r.Rate,
			Burst:  r.Burst,
			Period: time.Duration(r.Period) * time.Minute,
		})
		if res.Allowed <= 0 {
			return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded")
		}
		return c.Next()
	}
}

func (r *RedisRateLimiter) ThrottleByKeyAndIP(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		ip := c.IP()
		compositeKey := fmt.Sprintf(RateRequest, ip+"_"+key)
		res, _ := r.Limiter.Allow(ctx, compositeKey, redisrate.Limit{
			Rate:   r.Rate,
			Burst:  r.Burst,
			Period: time.Duration(r.Period) * time.Minute,
		})
		if res.Allowed <= 0 {
			return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded")
		}
		return c.Next()
	}
}
