package controllers

import (
	"ims/app/facades"
	"ims/app/models"

	"github.com/goravel/framework/contracts/http"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// ShowLogin displays login view
func (r *AuthController) ShowLogin(ctx http.Context) http.Response {
	return ctx.Response().View().Make("auth/login.tmpl", map[string]any{})
}

// Login processes user login
func (r *AuthController) Login(ctx http.Context) http.Response {
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	if email == "" || password == "" {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"Error": "Email and password are required.",
		})
	}

	var user models.User
	err := facades.Orm().Query().Where("email = ?", email).First(&user)
	if err != nil || user.ID == 0 {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"Error": "Invalid credentials.",
		})
	}

	if !facades.Hash().Check(password, user.Password) {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"Error": "Invalid credentials.",
		})
	}

	// Store in session
	ctx.Request().Session().Put("user_id", user.ID)

	// Record login activity
	RecordActivity(ctx, &user, "Login")

	if user.Role == "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/admin")
	}
	return ctx.Response().Redirect(http.StatusFound, "/orders")
}

// ShowRegister displays registration view
func (r *AuthController) ShowRegister(ctx http.Context) http.Response {
	return ctx.Response().View().Make("auth/register.tmpl", map[string]any{})
}

// Register processes user registration
func (r *AuthController) Register(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	if name == "" || email == "" || password == "" {
		return ctx.Response().View().Make("auth/register.tmpl", map[string]any{
			"Error": "All fields are required.",
		})
	}

	// Check if email exists
	count, err := facades.Orm().Query().Model(&models.User{}).Where("email = ?", email).Count()
	if err != nil {
		return ctx.Response().View().Make("auth/register.tmpl", map[string]any{
			"Error": "Registration error, try again.",
		})
	}
	if count > 0 {
		return ctx.Response().View().Make("auth/register.tmpl", map[string]any{
			"Error": "Email already registered.",
		})
	}

	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		return ctx.Response().View().Make("auth/register.tmpl", map[string]any{
			"Error": "Registration error, try again.",
		})
	}

	newUser := models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     "client",
	}

	if err := facades.Orm().Query().Create(&newUser); err != nil {
		return ctx.Response().View().Make("auth/register.tmpl", map[string]any{
			"Error": "Registration error, try again.",
		})
	}

	return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
		"Success": "Registration successful! Please login.",
	})
}

// Logout clears session and redirects
func (r *AuthController) Logout(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user != nil {
		RecordActivity(ctx, user, "Logout")
	}
	ctx.Request().Session().Forget("user_id")
	return ctx.Response().Redirect(http.StatusFound, "/")
}
