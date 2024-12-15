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
	courseRoute.Get("/:courseID", c.CourseController.Get)
	courseRoute.Put("/:courseID", c.CourseController.Update)

	// transactionRoute := c.App.Group("/api/transaction", c.EmployeeMiddleware)
	// transactionRoute.Post("/buy", c.TransactionController.Buy)
	// transactionRoute.Get("/:trxId", c.TransactionController.GetTransaction)

}
