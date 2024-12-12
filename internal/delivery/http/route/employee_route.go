package route

func (c *RouteConfig) SetupEmployeeRoute() {
	employeeRoute := c.App.Group("/api/employee", c.EmployeeAuthMiddleware)

	// employeeRoute.Patch("/update", c.EmployeeController.Update)
	employeeRoute.Post("/register", c.EmployeeController.Register)
	employeeRoute.Get("/current", c.EmployeeController.Current)
	employeeRoute.Delete("/logout", c.EmployeeController.Logout)

	courseCategoryRoute := employeeRoute.Group("/course-categories", c.EmployeeAuthMiddleware)

	courseCategoryRoute.Get("", c.CourseCatController.List)
	courseCategoryRoute.Post("", c.CourseCatController.Create)
	courseCategoryRoute.Get("/:courseCatID", c.CourseCatController.Get)
	courseCategoryRoute.Put("/:courseCatID", c.CourseCatController.Update)

	courseRoute := employeeRoute.Group("/course", c.EmployeeAuthMiddleware)

	courseRoute.Get("", c.CourseController.List)
	courseRoute.Post("", c.CourseController.Create)
	courseRoute.Get("/:courseID", c.CourseController.Get)
	courseRoute.Put("/:courseID", c.CourseController.Update)

	courseSectionRoute := courseRoute.Group("/section", c.EmployeeAuthMiddleware)

	courseSectionRoute.Post("", c.CourseSecController.Create)

	// transactionRoute := c.App.Group("/api/transaction", c.EmployeeMiddleware)
	// transactionRoute.Post("/buy", c.TransactionController.Buy)
	// transactionRoute.Get("/:trxId", c.TransactionController.GetTransaction)

}
