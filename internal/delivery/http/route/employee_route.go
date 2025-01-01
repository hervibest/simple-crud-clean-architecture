package route

func (c *RouteConfig) SetupEmployeeRoute() {
	employeeRoute := c.App.Group("/api/employee", c.EmployeeAuthMiddleware)

	// employeeRoute.Patch("/update", c.EmployeeController.Update)
	// employeeRoute.Post("/register", c.EmployeeController.Register)
	employeeRoute.Get("/current", c.EmployeeController.Current)
	employeeRoute.Delete("/logout", c.EmployeeController.Logout)

	courseCategoryRoute := employeeRoute.Group("/course-categories", c.EmployeeAuthMiddleware)

	courseCategoryRoute.Get("", c.CourseCatController.List)
	courseCategoryRoute.Post("", c.CourseCatController.Create)
	courseCategoryRoute.Get("/:courseCatID", c.CourseCatController.Get)
	courseCategoryRoute.Put("/:courseCatID", c.CourseCatController.Update)

	courseRoute := employeeRoute.Group("/course", c.EmployeeAuthMiddleware)

	courseSectionRoute := courseRoute.Group("/section", c.EmployeeAuthMiddleware)

	sectionVideoRoute := courseSectionRoute.Group("/video", c.EmployeeAuthMiddleware)

	sectionVideoRoute.Get("", c.SecVideoController.List)
	sectionVideoRoute.Get("/:secVideoID", c.SecVideoController.Get)
	sectionVideoRoute.Post("", c.SecVideoController.Create)
	sectionVideoRoute.Post("/upload/:secVideoID", c.SecVideoController.UploadVideo)
	sectionVideoRoute.Put("/:secVideoID", c.SecVideoController.Update)
	sectionVideoRoute.Delete("/:secVideoID", c.SecVideoController.Delete)

	courseSectionRoute.Get("", c.CourseSecController.List)
	courseSectionRoute.Get("/:courseSecId", c.CourseSecController.Get)
	courseSectionRoute.Post("", c.CourseSecController.Create)
	courseSectionRoute.Put("/:courseSecId", c.CourseSecController.Update)
	courseSectionRoute.Delete("/:courseSecId", c.CourseSecController.Delete)

	courseRoute.Get("", c.CourseController.List)
	courseRoute.Post("", c.CourseController.Create)
	courseRoute.Post("/upload/:courseId", c.CourseController.UploadThumbnail)
	courseRoute.Get("/:courseId", c.CourseController.Get)
	courseRoute.Put("/:courseId", c.CourseController.Update)

	discountRoute := employeeRoute.Group("/discount", c.EmployeeAuthMiddleware)

	discountRoute.Get("", c.DiscountController.List)
	discountRoute.Post("/create", c.DiscountController.Create)
	discountRoute.Get("/detail/:discountID", c.DiscountController.Get)
	discountRoute.Put("/update/:discountID", c.DiscountController.Update)

	voucherRoute := employeeRoute.Group("/voucher", c.EmployeeAuthMiddleware)

	voucherRoute.Get("", c.VoucherController.List)
	voucherRoute.Post("/create", c.VoucherController.Create)
	voucherRoute.Get("/detail/:voucherID", c.VoucherController.Get)
	voucherRoute.Put("/update/:voucherID", c.VoucherController.Update)
	voucherRoute.Post("/apply/:voucherCode", c.VoucherController.ApplyVoucher)

	transactionRoute := employeeRoute.Group("/transactions", c.EmployeeAuthMiddleware)

	transactionRoute.Get("", c.TransactionController.List)

	// transactionRoute := c.App.Group("/api/transaction", c.EmployeeMiddleware)
	// transactionRoute.Post("/buy", c.TransactionController.Buy)
	// transactionRoute.Get("/:trxId", c.TransactionController.GetTransaction)

	careerCategoryRoute := employeeRoute.Group("/career-categories", c.EmployeeAuthMiddleware)

	careerCategoryRoute.Get("", c.CourseCatController.List)
	careerCategoryRoute.Post("", c.CourseCatController.Create)
	careerCategoryRoute.Get("/:careerCatID", c.CourseCatController.Get)
	careerCategoryRoute.Put("/:careerCatID", c.CourseCatController.Update)

	careerRoute := employeeRoute.Group("/career", c.EmployeeAuthMiddleware)

	careerRoute.Get("", c.CourseController.List)
	careerRoute.Post("", c.CourseController.Create)
	careerRoute.Post("/upload/:careerId", c.CourseController.UploadThumbnail)
	careerRoute.Get("/:careerId", c.CourseController.Get)
	careerRoute.Put("/:careerId", c.CourseController.Update)

	//Certificate
	certificateCategoryRoute := employeeRoute.Group("/certificate-categories", c.EmployeeAuthMiddleware)
	certificateCategoryRoute.Get("", c.CertificateCatController.List)
	certificateCategoryRoute.Post("", c.CertificateCatController.Create)
	certificateCategoryRoute.Get("/:certificateCatID", c.CertificateCatController.Get)
	certificateCategoryRoute.Put("/:certificateCatID", c.CertificateCatController.Update)

	certificateRoute := employeeRoute.Group("/certificate", c.EmployeeAuthMiddleware)

	materialRoute := certificateRoute.Group("/material", c.EmployeeAuthMiddleware)
	materialRoute.Get("", c.CertifMaterialController.List)
	materialRoute.Post("/create", c.CertifMaterialController.Create)
	// materialRoute.Post("/upload/:careerId", c.CertifMaterialController.UploadThumbnail)
	materialRoute.Get("/detail/:materialID", c.CertifMaterialController.Get)
	materialRoute.Put("/update/:materialID", c.CertifMaterialController.Update)

	skkniRoute := certificateRoute.Group("/skkni", c.EmployeeAuthMiddleware)
	skkniRoute.Get("", c.CertifSkkniController.List)
	skkniRoute.Post("/create", c.CertifSkkniController.Create)
	skkniRoute.Post("/upload-thumbnail/:skkniId", c.CertifSkkniController.UploadThumbnail)
	skkniRoute.Get("/detail/:skkniId", c.CertifSkkniController.Get)
	skkniRoute.Put("/update/:skkniId", c.CertifSkkniController.Update)

	certificateRoute.Get("", c.CertificateController.List)
	certificateRoute.Post("/create", c.CertificateController.Create)
	certificateRoute.Post("/upload/:careerId", c.CertificateController.UploadThumbnail)
	certificateRoute.Get("/:careerId", c.CertificateController.Get)
	certificateRoute.Put("/:careerId", c.CertificateController.Update)

}
