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
	profileController := controllers.NewProfileController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		router.Get("/logout", authController.Logout)
		router.Get("/checkout/{product_id}", orderController.Checkout)
		router.Post("/checkout/{product_id}", orderController.PlaceOrder)
		router.Get("/orders", orderController.Index)
		
		// Profile & Account Settings
		router.Get("/profile", profileController.ShowProfile)
		router.Post("/profile", profileController.UpdateProfile)
		router.Get("/profile/password", profileController.ShowPassword)
		router.Post("/profile/password", profileController.UpdatePassword)
		router.Get("/profile/activity", profileController.ShowActivity)
	})

	// Admin Routes
	adminController := controllers.NewAdminController()
	facades.Route().Middleware(middleware.Admin()).Group(func(router route.Router) {
		router.Get("/admin", adminController.Dashboard)
		router.Post("/admin/orders/{order_id}/status", adminController.UpdateStatus)

		// Admin User Management
		router.Get("/admin/users", adminController.Users)
		router.Post("/admin/users", adminController.CreateUser)
		router.Post("/admin/users/{user_id}/delete", adminController.DeleteUser)

		// Admin CMS Management
		router.Get("/admin/cms", adminController.Cms)
		router.Post("/admin/cms/pages", adminController.CreateCmsPage)
		router.Post("/admin/cms/pages/{page_id}/edit", adminController.UpdateCmsPage)
		router.Post("/admin/cms/pages/{page_id}/delete", adminController.DeleteCmsPage)
		
		router.Get("/admin/cms/pages/{page_id}", adminController.CmsPageDetail)
		router.Post("/admin/cms/pages/{page_id}/contents", adminController.UpdateCmsPageContents)
		router.Post("/admin/cms/pages/{page_id}/contents/create", adminController.CreateCmsPageContent)

		// Admin Role & DB Schema Management
		router.Get("/admin/roles", adminController.Roles)
		router.Post("/admin/roles", adminController.CreateRole)
		router.Post("/admin/roles/{role_id}/edit", adminController.UpdateRole)
		router.Post("/admin/roles/{role_id}/delete", adminController.DeleteRole)

		// Database Schema Helpers
		router.Get("/admin/db-tables", adminController.GetDbTables)
		router.Get("/admin/db-columns/{table}", adminController.GetDbColumns)
	})
}
