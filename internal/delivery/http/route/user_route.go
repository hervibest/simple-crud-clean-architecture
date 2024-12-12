package route

func (c *RouteConfig) SetupUserRoute() {
	userRoute := c.App.Group("/api/users", c.UserAuthMiddleware)

	userRoute.Patch("/update", c.UserController.Update)
	userRoute.Get("/current", c.UserController.Current)
	userRoute.Delete("/logout", c.UserController.Logout)

	transactionRoute := c.App.Group("/api/transaction", c.UserAuthMiddleware)
	transactionRoute.Post("/buy", c.TransactionController.Buy)
	transactionRoute.Get("/:trxId", c.TransactionController.GetTransaction)

}
