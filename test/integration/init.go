package integration

import (
	"simple-crud-clean-architecture/internal/config"
	"simple-crud-clean-architecture/internal/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var app *fiber.App

var db *gorm.DB

var viperConfig *viper.Viper

var log *logrus.Logger

var redisClient *redis.Client

var tokenHelper *helper.TokenHelper

var emailHelper *helper.GomailSender

var midtrans *helper.MidtransClient

var customValidator helper.CustomValidator

var minioClient *helper.Minio

var job config.Job

func init() {
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	db = config.NewDatabase(viperConfig, log)
	app = config.NewFiber(viperConfig)
	redisClient = config.NewRedisClient(viperConfig, log)

	tokenHelper = helper.NewTokenHelper(viperConfig, log)
	emailHelper = helper.NewGomailSender(viperConfig, log)
	midtrans := helper.NewMidtransClient(viperConfig, log)
	customValidator = helper.NewCustomValidator(viperConfig)
	minioClient = helper.NewMinio(viperConfig, log)
	job = config.NewCronJob(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:              db,
		App:             app,
		Log:             log,
		Config:          viperConfig,
		Redis:           redisClient,
		TokenHelper:     tokenHelper,
		EmailHelper:     emailHelper,
		Midtrans:        midtrans,
		MinioClient:     minioClient,
		CustomValidator: customValidator,
		Job:             job,
	})

	webPort := viperConfig.GetInt("web.port")
	env := viperConfig.GetString("app.env")

	if env == "development" {
		helper.StartServer(app, webPort)
	} else if env == "production" {
		helper.StartServerWithGracefulShutdown(app, webPort)
	}
}
