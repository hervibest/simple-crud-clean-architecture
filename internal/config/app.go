package config

import (
	http "simple-crud-clean-architecture/internal/delivery/http/controller"
	"simple-crud-clean-architecture/internal/delivery/http/middleware"
	"simple-crud-clean-architecture/internal/delivery/http/route"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/repository"
	"simple-crud-clean-architecture/internal/usecase"

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
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, userRepository, config.Redis, config.TokenHelper, config.EmailHelper, config.CustomValidator)
	employeeUseCase := usecase.NewEmployeeUseCase(config.DB, config.Log, employeeRepository, config.Redis, config.TokenHelper, config.EmailHelper, config.CustomValidator)
	courseCatUseCase := usecase.NewCourseCatUseCase(config.DB, config.Log, courseCatRepository, config.CustomValidator)
	courseUseCase := usecase.NewCourseUseCase(config.DB, config.Log, courseRepository, courseCatRepository, config.MinioClient, config.CustomValidator)
	courseSecUseCase := usecase.NewCourseSecUseCase(config.DB, config.Log, courseSecRepository, courseRepository, config.CustomValidator)
	courseVidUseCase := usecase.NewSecVideoUseCase(config.DB, config.Log, courseSecRepository, courseVidRepository, config.MinioClient, config.CustomValidator)
	careerCatUseCase := usecase.NewCareerCatUseCase(config.DB, config.Log, careerCatRepository, config.CustomValidator)
	careerUseCase := usecase.NewCareerUseCase(config.DB, config.Log, careerRepository, careerCatRepository, config.MinioClient, config.CustomValidator)

	discountUseCase := usecase.NewDiscountUseCase(config.DB, config.Log, discountRepository, courseRepository, config.CustomValidator)
	voucherUseCase := usecase.NewVoucherUseCase(config.DB, config.Log, voucherRepository, courseRepository, config.CustomValidator)
	transactionUseCase := usecase.NewTransactionUseCase(config.DB, config.Log, transactionRepository, courseRepository, userRepository, voucherRepository, config.Midtrans, config.CustomValidator)

	certifCatUseCase := usecase.NewCertifCatUseCase(config.DB, config.Log, certifCatRepository, config.CustomValidator)
	certifUseCase := usecase.NewCertificateUseCase(config.DB, config.Log, certifiRepository, certifCatRepository, config.MinioClient, config.CustomValidator)
	certifMaterialUseCase := usecase.NewCertifMaterialUseCase(config.DB, config.Log, certifMaterialRepository, certifiRepository, config.CustomValidator)
	certifSkkniUseCase := usecase.NewCertifSkkniUseCase(config.DB, config.Log, certifSkkniRepository, certifiRepository, config.MinioClient, config.CustomValidator)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	employeeController := http.NewEmployeeController(employeeUseCase, config.Log)

	courseCatController := http.NewCourseCatController(courseCatUseCase, config.Log)
	courseController := http.NewCourseController(courseUseCase, config.Log)
	courseSecController := http.NewCourseSecController(courseSecUseCase, config.Log)
	courseVidControler := http.NewSecVideoController(courseVidUseCase, config.Log, config.MinioClient)

	careerCatController := http.NewCareerCatController(careerCatUseCase, config.Log)
	careerController := http.NewCareerController(careerUseCase, config.Log)

	discountController := http.NewDiscountController(discountUseCase, config.Log)
	voucherController := http.NewVoucherController(voucherUseCase, config.Log)
	transactionController := http.NewTransactionController(transactionUseCase, config.Log, config.Midtrans)

	certifCatController := http.NewCertificateCatController(certifCatUseCase, config.Log)
	certifController := http.NewCertificateController(certifUseCase, config.Log)
	certifMatController := http.NewCertifMaterialController(certifMaterialUseCase, config.Log)
	certifSkkniController := http.NewCertifSkkniController(certifSkkniUseCase, config.Log)

	// upload controller
	uploadController := http.NewUploadController(config.Log, *config.MinioClient)

	// setup throttle
	throttle := middleware.NewRedisRateLimiter(config.Redis, config.Config)

	//setup
	userAuthMiddleware := middleware.NewUserAuth(userUseCase, config.CustomValidator)
	buyableCourseMiddleware := middleware.NewBuyableCourse(courseUseCase, config.CustomValidator)
	employeeAuthMiddleware := middleware.NewEmployeeAuth(employeeUseCase, config.CustomValidator)

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
