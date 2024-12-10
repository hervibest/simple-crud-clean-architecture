package route

import (
	http "simple-crud-clean-architecture/internal/delivery/http/controller"
	"simple-crud-clean-architecture/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController

	Throttle       *middleware.Throttle
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupUserRoute()
}
