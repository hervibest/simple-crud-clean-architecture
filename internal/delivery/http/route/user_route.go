package route

func (c *RouteConfig) SetupUserRoute() {
	userRoute := c.App.Group("/api/users", c.AuthMiddleware)

	userRoute.Patch("/update", c.UserController.Update)
	userRoute.Get("/current", c.UserController.Current)
	userRoute.Delete("/logout", c.UserController.Logout)

	courseCategoryRoute := c.App.Group("/api/course-categories", c.AuthMiddleware)

	courseCategoryRoute.Get("", c.CourseCatController.List)
	courseCategoryRoute.Post("", c.CourseCatController.Create)
	courseCategoryRoute.Get("/:courseCatID", c.CourseCatController.Get)
	courseCategoryRoute.Put("/:courseCatID", c.CourseCatController.Update)

	courseRoute := c.App.Group("/api/courses", c.AuthMiddleware)

	courseRoute.Get("", c.CourseController.List)
	courseRoute.Post("", c.CourseController.Create)
	courseRoute.Get("/:courseID", c.CourseController.Get)
	courseRoute.Put("/:courseID", c.CourseController.Update)

}
