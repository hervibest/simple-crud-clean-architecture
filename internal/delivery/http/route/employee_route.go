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

	// transactionRoute := c.App.Group("/api/transaction", c.EmployeeMiddleware)
	// transactionRoute.Post("/buy", c.TransactionController.Buy)
	// transactionRoute.Get("/:trxId", c.TransactionController.GetTransaction)

}
