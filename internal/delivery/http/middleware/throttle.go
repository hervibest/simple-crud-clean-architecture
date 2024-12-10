package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// Throttle struct to hold configuration
type Throttle struct {
	max int
	sec int
}

// NewThrottle creates a new Throttle instance
func NewThrottle(max, sec int) *Throttle {
	return &Throttle{
		max: max,
		sec: sec,
	}
}

// ThrottleByKey method for Throttle struct
func (t *Throttle) ThrottleByKey(key string) func(c *fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        t.max,
		Expiration: time.Duration(t.sec) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return key
		},
		LimitReached: func(c *fiber.Ctx) error {
			fmt.Println("LimitReached")
			return fiber.ErrTooManyRequests
		},
	})
}

// ThrottleByKeyAndIP method for Throttle struct
func (t *Throttle) ThrottleByKeyAndIP(key string) func(c *fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        t.max,
		Expiration: time.Duration(t.sec) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return string(c.IP()) + key
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	})
}

// ThrottleByIp method for Throttle struct
func (t *Throttle) ThrottleByIp(c *fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        t.max,
		Expiration: time.Duration(t.sec) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	})(c)
}
