package route

func (c *RouteConfig) SetupGuestRoute() {
	// GUEST USER ROUTE
	userRoutes := c.App.Group("/api/user")
	userRoutes.Post("/register", c.UserController.Register)
	userRoutes.Post("/login", c.UserController.Login)
	userRoutes.Post("/verify/:token", c.UserController.VerifyEmail)
	userRoutes.Post("/request-resend-email", c.Throttle.ThrottleByKey("request-resend-email"), c.UserController.RequestEmailVerification)
	userRoutes.Post("/request-access-token", c.UserController.RequestAccessToken)

	userRoutes.Post("/reset-password/request", c.Throttle.ThrottleByKey("request-resend-password"), c.UserController.RequestResetPassword)
	userRoutes.Post("/reset-password/validate", c.UserController.ValidateResetToken)
	userRoutes.Post("/reset-password/reset/:token", c.UserController.ResetPassword)

	c.App.Post("/api/admin/login", c.EmployeeController.Login)

	c.App.Post("/api/webhook/notify", c.TransactionController.Notify)

	// c.App.Post("api/upload", c.UploadController.UploadFile)
	// c.App.Delete("api/upload/:fileName", c.UploadController.DeleteFile)

	// c.App.Post("/api/register", c.EmployeeController.Register)

}
