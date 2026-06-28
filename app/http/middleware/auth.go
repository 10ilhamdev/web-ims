package middleware

import (
	"ims/app/facades"
	"ims/app/models"

	"github.com/goravel/framework/contracts/http"
)

// Auth checks if user is logged in
func Auth() http.Middleware {
	return func(ctx http.Context) {
		userID := ctx.Request().Session().Get("user_id")
		if userID == nil {
			ctx.Response().Redirect(http.StatusFound, "/login")
			return
		}
		ctx.Request().Next()
	}
}

// Guest redirects logged in users to dashboard/home
func Guest() http.Middleware {
	return func(ctx http.Context) {
		userID := ctx.Request().Session().Get("user_id")
		if userID != nil {
			ctx.Response().Redirect(http.StatusFound, "/")
			return
		}
		ctx.Request().Next()
	}
}

// Admin checks if logged in user is an admin
func Admin() http.Middleware {
	return func(ctx http.Context) {
		userID := ctx.Request().Session().Get("user_id")
		if userID == nil {
			ctx.Response().Redirect(http.StatusFound, "/login")
			return
		}

		var user models.User
		err := facades.Orm().Query().Find(&user, userID)
		if err != nil || user.ID == 0 || user.Role != "admin" {
			ctx.Response().Redirect(http.StatusFound, "/")
			return
		}

		ctx.Request().Next()
	}
}
