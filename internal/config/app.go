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
	DB              *gorm.DB
	App             *fiber.App
	Log             *logrus.Logger
	Validate        *validator.Validate
	Config          *viper.Viper
	Redis           *redis.Client
	TokenHelper     *helper.TokenHelper
	EmailHelper     *helper.GomailSender
	Midtrans        *helper.MidtransClient
	MinioClient     *helper.Minio
	CustomValidator helper.CustomValidator
	Job             Job
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	employeeRepository := repository.NewEmployeeRepository(config.Log)

	courseCatRepository := repository.NewCourseCatRepository(config.Log)
	courseRepository := repository.NewCourseRepository(config.Log)
	courseSecRepository := repository.NewCourseSectionRepository(config.Log)
	courseVidRepository := repository.NewSecVideoRepository(config.Log)

	careerCatRepository := repository.NewCareerCatRepository(config.Log)
	careerRepository := repository.NewCareerRepository(config.Log)

	discountRepository := repository.NewDiscountRepository(config.Log)
	voucherRepository := repository.NewVoucherRepository(config.Log)
	transactionRepository := repository.NewTransactionRepository(config.Log)

	certifCatRepository := repository.NewCertifCategoryRepository(config.Log)
	certifiRepository := repository.NewCertificateRepository(config.Log)
	certifMaterialRepository := repository.NewCertifMaterialRepository(config.Log)
	certifSkkniRepository := repository.NewCertifSkkniRepository(config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, config.Redis, config.TokenHelper, config.EmailHelper)
	employeeUseCase := usecase.NewEmployeeUseCase(config.DB, config.Log, config.Validate, employeeRepository, config.Redis, config.TokenHelper, config.EmailHelper)
	courseCatUseCase := usecase.NewCourseCatUseCase(config.DB, config.Log, config.Validate, courseCatRepository)
	courseUseCase := usecase.NewCourseUseCase(config.DB, config.Log, config.Validate, courseRepository, courseCatRepository, config.MinioClient)
	courseSecUseCase := usecase.NewCourseSecUseCase(config.DB, config.Log, config.Validate, courseSecRepository, courseRepository)
	courseVidUseCase := usecase.NewSecVideoUseCase(config.DB, config.Log, config.Validate, courseSecRepository, courseVidRepository, config.MinioClient)
	careerCatUseCase := usecase.NewCareerCatUseCase(config.DB, config.Log, config.Validate, careerCatRepository)
	careerUseCase := usecase.NewCareerUseCase(config.DB, config.Log, config.Validate, careerRepository, careerCatRepository, config.MinioClient)

	discountUseCase := usecase.NewDiscountUseCase(config.DB, config.Log, discountRepository, courseRepository)
	voucherUseCase := usecase.NewVoucherUseCase(config.DB, config.Log, voucherRepository, courseRepository)
	transactionUseCase := usecase.NewTransactionUseCase(config.DB, config.Log, config.Validate, transactionRepository, courseRepository, userRepository, voucherRepository, config.Midtrans)

	certifCatUseCase := usecase.NewCertifCatUseCase(config.DB, config.Log, config.Validate, certifCatRepository)
	certifUseCase := usecase.NewCertificateUseCase(config.DB, config.Log, config.Validate, certifiRepository, certifCatRepository, config.MinioClient)
	certifMaterialUseCase := usecase.NewCertifMaterialUseCase(config.DB, config.Log, config.Validate, certifMaterialRepository, certifiRepository)
	certifSkkniUseCase := usecase.NewCertifSkkniUseCase(config.DB, config.Log, config.Validate, certifSkkniRepository, certifiRepository, config.MinioClient)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log, config.CustomValidator)
	employeeController := http.NewEmployeeController(employeeUseCase, config.Log, config.CustomValidator)

	courseCatController := http.NewCourseCatController(courseCatUseCase, config.Log, config.CustomValidator)
	courseController := http.NewCourseController(courseUseCase, config.Log, config.CustomValidator)
	courseSecController := http.NewCourseSecController(courseSecUseCase, config.Log, config.CustomValidator)
	courseVidControler := http.NewSecVideoController(courseVidUseCase, config.Log, config.CustomValidator, config.MinioClient)

	careerCatController := http.NewCareerCatController(careerCatUseCase, config.Log, config.CustomValidator)
	careerController := http.NewCareerController(careerUseCase, config.Log, config.CustomValidator)

	discountController := http.NewDiscountController(discountUseCase, config.Log, config.CustomValidator)
	voucherController := http.NewVoucherController(voucherUseCase, config.Log, config.CustomValidator)
	transactionController := http.NewTransactionController(transactionUseCase, config.Log, config.Midtrans, config.CustomValidator)

	certifCatController := http.NewCertificateCatController(certifCatUseCase, config.Log, config.CustomValidator)
	certifController := http.NewCertificateController(certifUseCase, config.Log, config.CustomValidator)
	certifMatController := http.NewCertifMaterialController(certifMaterialUseCase, config.Log, config.CustomValidator)
	certifSkkniController := http.NewCertifSkkniController(certifSkkniUseCase, config.Log, config.CustomValidator)

	// upload controller
	uploadController := http.NewUploadController(config.Log, config.CustomValidator, *config.MinioClient)

	// setup throttle
	throttle := middleware.NewRedisRateLimiter(config.Redis, config.Config)

	//setup
	userAuthMiddleware := middleware.NewUserAuth(userUseCase)
	buyableCourseMiddleware := middleware.NewBuyableCourse(courseUseCase)
	employeeAuthMiddleware := middleware.NewEmployeeAuth(employeeUseCase)

	// Scheduler
	job := config.Job
	job.RunCron(discountUseCase)

	routeConfig := route.RouteConfig{
		App: config.App,

		UserController:     userController,
		EmployeeController: employeeController,

		CourseCatController: courseCatController,
		CourseController:    courseController,
		CourseSecController: courseSecController,
		SecVideoController:  courseVidControler,

		CareerCatController: careerCatController,
		CareerController:    careerController,

		DiscountController:    discountController,
		VoucherController:     voucherController,
		UploadController:      uploadController,
		TransactionController: transactionController,

		CertifMaterialController: certifMatController,
		CertifSkkniController:    certifSkkniController,
		CertificateCatController: certifCatController,
		CertificateController:    certifController,

		Throttle:                throttle,
		UserAuthMiddleware:      userAuthMiddleware,
		BuyableCourseMiddleware: buyableCourseMiddleware,
		EmployeeAuthMiddleware:  employeeAuthMiddleware,
	}
	routeConfig.Setup()
}
