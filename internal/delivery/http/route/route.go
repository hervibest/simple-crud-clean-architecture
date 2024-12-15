package route

import (
	http "simple-crud-clean-architecture/internal/delivery/http/controller"
	"simple-crud-clean-architecture/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                 *fiber.App
	UserController      *http.UserController
	EmployeeController  *http.EmployeeController
	CourseCatController *http.CourseCatController
	CourseController    *http.CourseController
	CourseSecController *http.CourseSectionController
	SecVideoController  *http.SectionVideoController

	DiscountController    *http.DiscountController
	TransactionController *http.TransactionController

	Throttle                *middleware.Throttle
	UserAuthMiddleware      fiber.Handler
	BuyableCourseMiddleware fiber.Handler
	EmployeeAuthMiddleware  fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupUserRoute()
	c.SetupEmployeeRoute()
}
