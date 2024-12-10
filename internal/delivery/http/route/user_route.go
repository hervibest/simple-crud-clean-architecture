package route

func (c *RouteConfig) SetupUserRoute() {
	userRoute := c.App.Group("/api/users", c.AuthMiddleware)

	userRoute.Patch("/update", c.UserController.Update)
	userRoute.Get("/current", c.UserController.Current)
	userRoute.Delete("/logout", c.UserController.Logout)
}
