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

	// Chat API Routes
	chatController := controllers.NewChatController()
	facades.Route().Post("/api/chat/ai", chatController.AskAI)
	facades.Route().Post("/api/chat/support/init", chatController.InitSupport)
	facades.Route().Post("/api/chat/support/send", chatController.SendSupportMessage)
	facades.Route().Get("/api/chat/support/messages", chatController.GetSupportMessages)

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

		// Admin Support Chat
		router.Get("/admin/chat", chatController.AdminChatView)
		router.Get("/api/admin/chat/sessions", chatController.AdminGetSessions)
		router.Post("/api/admin/chat/reply", chatController.AdminSendReply)

		// Admin User Management
		router.Get("/admin/users", adminController.Users)
		router.Get("/admin/users/create", adminController.CreateUserForm)
		router.Post("/admin/users", adminController.CreateUser)
		router.Get("/admin/users/{user_id}", adminController.UserDetail)
		router.Get("/admin/users/{user_id}/edit", adminController.EditUserForm)
		router.Post("/admin/users/{user_id}/edit", adminController.UpdateUser)
		router.Post("/admin/users/{user_id}/delete", adminController.DeleteUser)

		// Admin CMS Management
		router.Get("/admin/cms", adminController.Cms)
		router.Post("/admin/cms/pages", adminController.CreateCmsPage)
		router.Post("/admin/cms/pages/{page_id}/edit", adminController.UpdateCmsPage)
		router.Post("/admin/cms/pages/{page_id}/delete", adminController.DeleteCmsPage)
		
		router.Get("/admin/cms/pages/{page_id}", adminController.CmsPageDetail)
		router.Post("/admin/cms/pages/{page_id}/contents", adminController.UpdateCmsPageContents)
		router.Post("/admin/cms/pages/{page_id}/contents/create", adminController.CreateCmsPageContent)
		router.Post("/admin/cms/products/{product_id}/price", adminController.UpdateProductPrice)
		router.Post("/admin/cms/contents/{content_id}/delete", adminController.DeleteCmsPageContent)

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
