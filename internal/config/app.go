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
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	courseCatRepository := repository.NewCourseCatRepository(config.Log)
	courseRepositoru := repository.NewCourseRepository(config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, config.Redis, config.TokenHelper, config.EmailHelper)
	courseCatUseCase := usecase.NewCourseCatUseCase(config.DB, config.Log, config.Validate, courseCatRepository)
	courseUseCase := usecase.NewCourseUseCase(config.DB, config.Log, config.Validate, courseRepositoru, courseCatRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	courseCatController := http.NewCourseCatController(courseCatUseCase, config.Log)
	courseController := http.NewCourseController(courseUseCase, config.Log)

	// setup throttle
	throttle := middleware.NewThrottle(1, 60)

	//setup
	middleware := middleware.NewAuth(userUseCase)

	routeConfig := route.RouteConfig{
		App:                 config.App,
		UserController:      userController,
		CourseCatController: courseCatController,
		CourseController:    courseController,
		Throttle:            throttle,
		AuthMiddleware:      middleware,
	}
	routeConfig.Setup()
}
