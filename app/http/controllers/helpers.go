package controllers

import (
	"ims/app/facades"
	"ims/app/models"

	"github.com/goravel/framework/contracts/http"
)

// GetCurrentUser returns the currently logged in user or nil
func GetCurrentUser(ctx http.Context) *models.User {
	userID := ctx.Request().Session().Get("user_id")
	if userID == nil {
		return nil
	}
	
	// Convert session value to uint safely
	var uID uint
	switch v := userID.(type) {
	case uint:
		uID = v
	case int:
		uID = uint(v)
	case float64:
		uID = uint(v)
	default:
		return nil
	}

	var user models.User
	err := facades.Orm().Query().Find(&user, uID)
	if err != nil || user.ID == 0 {
		return nil
	}
	return &user
}
