package routes

import (
	"github.com/goravel/framework/contracts/route"

	"ims/app/http/controllers"
	"ims/app/http/middleware"
	"ims/app/facades"
)

func Web() {
	facades.Route().Static("public", "./public")

	// Guest Routes
	guestController := controllers.NewGuestController()
	facades.Route().Get("/", guestController.Index)
	facades.Route().Get("/services", guestController.Services)
	facades.Route().Get("/products", guestController.Products)

	// Auth Guest Routes
	authController := controllers.NewAuthController()
	facades.Route().Middleware(middleware.Guest()).Group(func(router route.Router) {
		router.Get("/login", authController.ShowLogin)
		router.Post("/login", authController.Login)
		router.Get("/register", authController.ShowRegister)
		router.Post("/register", authController.Register)
	})

	// Logged-in Client Routes
	orderController := controllers.NewOrderController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		router.Get("/logout", authController.Logout)
		router.Get("/checkout/{product_id}", orderController.Checkout)
		router.Post("/checkout/{product_id}", orderController.PlaceOrder)
		router.Get("/orders", orderController.Index)
	})

	// Admin Routes
	adminController := controllers.NewAdminController()
	facades.Route().Middleware(middleware.Admin()).Group(func(router route.Router) {
		router.Get("/admin", adminController.Dashboard)
		router.Post("/admin/orders/{order_id}/status", adminController.UpdateStatus)
	})
}
