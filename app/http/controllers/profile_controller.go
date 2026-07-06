package controllers

import (
	"ims/app/facades"
	"ims/app/models"

	"github.com/goravel/framework/contracts/http"
)

type ProfileController struct{}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}

// ShowProfile renders user profile data read-only with edit toggle
func (r *ProfileController) ShowProfile(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	successMsg := ctx.Request().Session().Get("success")
	errorMsg := ctx.Request().Session().Get("error")
	ctx.Request().Session().Forget("success")
	ctx.Request().Session().Forget("error")

	return ctx.Response().View().Make("profile/info.tmpl", map[string]any{
		"User":    user,
		"Success": successMsg,
		"Error":   errorMsg,
	})
}

// UpdateProfile processes personal data updates
func (r *ProfileController) UpdateProfile(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")

	if name == "" || email == "" {
		ctx.Request().Session().Put("error", "Name and email are required.")
		return ctx.Response().Redirect(http.StatusFound, "/profile")
	}

	user.Name = name
	user.Email = email
	_ = facades.Orm().Query().Save(&user)

	RecordActivity(ctx, user, "Updated Profile Info")
	ctx.Request().Session().Put("success", "Profil berhasil diperbarui!")

	return ctx.Response().Redirect(http.StatusFound, "/profile")
}

// ShowPassword renders password update view
func (r *ProfileController) ShowPassword(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	successMsg := ctx.Request().Session().Get("success")
	errorMsg := ctx.Request().Session().Get("error")
	ctx.Request().Session().Forget("success")
	ctx.Request().Session().Forget("error")

	return ctx.Response().View().Make("profile/password.tmpl", map[string]any{
		"User":    user,
		"Success": successMsg,
		"Error":   errorMsg,
	})
}

// UpdatePassword updates user password
func (r *ProfileController) UpdatePassword(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	currentPassword := ctx.Request().Input("current_password")
	newPassword := ctx.Request().Input("new_password")

	if currentPassword == "" || newPassword == "" {
		ctx.Request().Session().Put("error", "Semua kolom password wajib diisi.")
		return ctx.Response().Redirect(http.StatusFound, "/profile/password")
	}

	if !facades.Hash().Check(currentPassword, user.Password) {
		ctx.Request().Session().Put("error", "Kata sandi lama tidak sesuai.")
		return ctx.Response().Redirect(http.StatusFound, "/profile/password")
	}

	hashedPassword, err := facades.Hash().Make(newPassword)
	if err != nil {
		ctx.Request().Session().Put("error", "Gagal memperbarui kata sandi.")
		return ctx.Response().Redirect(http.StatusFound, "/profile/password")
	}

	user.Password = hashedPassword
	_ = facades.Orm().Query().Save(&user)

	RecordActivity(ctx, user, "Changed Password")
	ctx.Request().Session().Put("success", "Kata sandi berhasil diperbarui!")

	return ctx.Response().Redirect(http.StatusFound, "/profile/password")
}

// ShowActivity lists user audit logs
func (r *ProfileController) ShowActivity(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var logs []models.ActivityLog
	err := facades.Orm().Query().Where("user_id = ?", user.ID).Order("created_at desc").Limit(30).Get(&logs)
	if err != nil {
		logs = []models.ActivityLog{}
	}

	return ctx.Response().View().Make("profile/activity.tmpl", map[string]any{
		"User": user,
		"Logs": logs,
	})
}
