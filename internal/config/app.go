package config

import (
	http "simple-crud-clean-architecture/internal/delivery/http/controller"
	"simple-crud-clean-architecture/internal/delivery/http/middleware"
	"simple-crud-clean-architecture/internal/delivery/http/route"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/repository"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Log         *logrus.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	Redis       *redis.Client
	TokenHelper *helper.TokenHelper
	EmailHelper *helper.GomailSender
	Midtrans    *helper.MidtransClient
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	employeeRepository := repository.NewEmployeeRepository(config.Log)

	courseCatRepository := repository.NewCourseCatRepository(config.Log)
	courseRepository := repository.NewCourseRepository(config.Log)
	transactionRepository := repository.NewTransactionRepository(config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, config.Redis, config.TokenHelper, config.EmailHelper)
	employeeUseCase := usecase.NewEmployeeUseCase(config.DB, config.Log, config.Validate, employeeRepository, config.Redis, config.TokenHelper, config.EmailHelper)
	courseCatUseCase := usecase.NewCourseCatUseCase(config.DB, config.Log, config.Validate, courseCatRepository)
	courseUseCase := usecase.NewCourseUseCase(config.DB, config.Log, config.Validate, courseRepository, courseCatRepository)
	transactionUseCase := usecase.NewTransactionUseCase(config.DB, config.Log, config.Validate, transactionRepository, courseRepository, userRepository, config.Midtrans)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	employeeController := http.NewEmployeeController(employeeUseCase, config.Log)

	courseCatController := http.NewCourseCatController(courseCatUseCase, config.Log)
	courseController := http.NewCourseController(courseUseCase, config.Log)
	transactionController := http.NewTransactionController(transactionUseCase, config.Log, config.Midtrans)

	// setup throttle
	throttle := middleware.NewThrottle(1, 60)

	//setup
	userAuthMiddleware := middleware.NewUserAuth(userUseCase)
	buyableCourseMiddleware := middleware.NewBuyableCourse(courseUseCase)
	employeeAuthMiddleware := middleware.NewEmployeeAuth(employeeUseCase)

	routeConfig := route.RouteConfig{
		App: config.App,

		UserController:     userController,
		EmployeeController: employeeController,

		CourseCatController:   courseCatController,
		CourseController:      courseController,
		TransactionController: transactionController,

		Throttle:                throttle,
		UserAuthMiddleware:      userAuthMiddleware,
		BuyableCourseMiddleware: buyableCourseMiddleware,
		EmployeeAuthMiddleware:  employeeAuthMiddleware,
	}
	routeConfig.Setup()
}
