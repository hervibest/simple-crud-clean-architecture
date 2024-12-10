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

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, config.Redis, config.TokenHelper, config.EmailHelper)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)

	// setup throttle
	throttle := middleware.NewThrottle(1, 60)

	//setup

	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
		Throttle:       throttle,
	}
	routeConfig.Setup()
}
