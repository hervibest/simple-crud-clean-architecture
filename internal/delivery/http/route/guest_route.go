package route

func (c *RouteConfig) SetupGuestRoute() {
	// GUEST USER ROUTE
	userRoutes := c.App.Group("/api/user")
	userRoutes.Post("/register", c.UserController.Register)
	userRoutes.Post("/verify/:token", c.UserController.VerifyEmail)

	userRoutes.Post("/request-resend-email", c.Throttle.ThrottleByKey("request-resend-email"), c.UserController.RequestEmailVerification)

}
