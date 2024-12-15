package main

import (
	"simple-crud-clean-architecture/internal/config"
	"simple-crud-clean-architecture/internal/helper"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	app := config.NewFiber(viperConfig)
	redis := config.NewRedisClient(viperConfig, log)
	validate := config.NewValidator(viperConfig)

	tokenHelper := helper.NewTokenHelper(viperConfig, log)
	emailHelper := helper.NewGomailSender(viperConfig, log)
	midtrans := helper.NewMidtransClient(viperConfig, log)
	customValidator := helper.NewCustomValidator(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:              db,
		App:             app,
		Log:             log,
		Validate:        validate,
		Config:          viperConfig,
		Redis:           redis,
		TokenHelper:     tokenHelper,
		EmailHelper:     emailHelper,
		Midtrans:        midtrans,
		CustomValidator: customValidator,
	})

	webPort := viperConfig.GetInt("web.port")
	env := viperConfig.GetString("app.env")

	if env == "development" {
		helper.StartServer(app, webPort)
	} else if env == "production" {
		helper.StartServerWithGracefulShutdown(app, webPort)
	}
}
