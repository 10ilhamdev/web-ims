package controllers

import (
	"encoding/json"
	"ims/app/facades"
	"ims/app/models"
	"strings"

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

// GetDeviceType parses user agent string to return device type
func GetDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipod") {
		return "Mobile"
	}
	if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		return "Tablet"
	}
	return "Desktop"
}

// RecordActivity stores user action in activity_logs table
func RecordActivity(ctx http.Context, user *models.User, activity string) {
	if user == nil {
		return
	}
	ua := ctx.Request().Header("User-Agent")
	device := GetDeviceType(ua)
	ip := ctx.Request().Ip()
	if ip == "" {
		ip = "127.0.0.1"
	}

	log := models.ActivityLog{
		UserID:   user.ID,
		Activity: activity,
		Device:   device,
		IP:       ip,
	}
	_ = facades.Orm().Query().Create(&log)
}

// GetCmsTranslations queries guest contents and returns ID and EN translation maps
func GetCmsTranslations(ctx http.Context) (map[string]string, map[string]string) {
	var contents []models.GuestContent
	_ = facades.Orm().Query().Get(&contents)

	idMap := make(map[string]string)
	enMap := make(map[string]string)
	for _, c := range contents {
		idMap[c.Key] = c.ValueId
		enMap[c.Key] = c.ValueEn
	}
	return idMap, enMap
}

// GetCmsTranslationsJson queries guest contents and returns ID and EN translation JSON strings
func GetCmsTranslationsJson(ctx http.Context) (string, string) {
	var contents []models.GuestContent
	_ = facades.Orm().Query().Get(&contents)

	idMap := make(map[string]string)
	enMap := make(map[string]string)
	for _, c := range contents {
		idMap[c.Key] = c.ValueId
		enMap[c.Key] = c.ValueEn
	}
	idBytes, _ := json.Marshal(idMap)
	enBytes, _ := json.Marshal(enMap)
	return string(idBytes), string(enBytes)
}
