
package controllers

import (
	"ims/app/facades"
	"ims/app/models"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
)

type AdminController struct{}

func NewAdminController() *AdminController {
	return &AdminController{}
}

// Dashboard serves AdminLTE style panel
func (r *AdminController) Dashboard(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var orders []models.Order
	err := facades.Orm().Query().Get(&orders)
	if err == nil {
		for i := range orders {
			var client models.User
			_ = facades.Orm().Query().Find(&client, orders[i].UserID)
			orders[i].User = &client

			var product models.Product
			_ = facades.Orm().Query().Find(&product, orders[i].ProductID)
			orders[i].Product = &product
		}
	} else {
		orders = []models.Order{}
	}

	// Calculate Stats
	var totalOrders = len(orders)
	var activeOrders = 0
	var totalRevenue float64 = 0
	var totalClients int64 = 0

	for _, o := range orders {
		if o.Status == "in_progress" || o.Status == "pending" {
			activeOrders++
		}
		if o.Status == "completed" {
			totalRevenue += o.Price
		}
	}

	totalClients, _ = facades.Orm().Query().Model(&models.User{}).Where("role = ?", "client").Count()

	return ctx.Response().View().Make("admin/dashboard.tmpl", map[string]any{
		"User":          user,
		"Orders":        orders,
		"TotalOrders":   totalOrders,
		"ActiveOrders":  activeOrders,
		"TotalRevenue":  totalRevenue,
		"TotalClients":  totalClients,
	})
}

// UpdateStatus changes order status
func (r *AdminController) UpdateStatus(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	oIDStr := ctx.Request().Route("order_id")
	oID, err := strconv.Atoi(oIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin")
	}

	status := ctx.Request().Input("status")

	var order models.Order
	err = facades.Orm().Query().Find(&order, oID)
	if err == nil && order.ID != 0 {
		order.Status = status
		_ = facades.Orm().Query().Save(&order)
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin")
}

// Users lists all users for admin management
func (r *AdminController) Users(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var users []models.User
	_ = facades.Orm().Query().Get(&users)

	return ctx.Response().View().Make("admin/users.tmpl", map[string]any{
		"User":  user,
		"Users": users,
	})
}

// CreateUserForm shows the create user page
func (r *AdminController) CreateUserForm(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	return ctx.Response().View().Make("admin/users_create.tmpl", map[string]any{
		"User": user,
	})
}

// CreateUser handles admin adding a new user
func (r *AdminController) CreateUser(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")
	role := ctx.Request().Input("role")

	if name == "" || email == "" || password == "" {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users/create")
	}

	hashedPassword, err := facades.Hash().Make(password)
	if err == nil {
		newUser := models.User{
			Name:     name,
			Email:    email,
			Password: hashedPassword,
			Role:     role,
		}
		_ = facades.Orm().Query().Create(&newUser)
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/users")
}

// UserDetail shows dynamic details about a specific user
func (r *AdminController) UserDetail(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	uIDStr := ctx.Request().Route("user_id")
	uID, err := strconv.Atoi(uIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users")
	}

	var targetUser models.User
	err = facades.Orm().Query().Find(&targetUser, uID)
	if err != nil || targetUser.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users")
	}

	return ctx.Response().View().Make("admin/users_detail.tmpl", map[string]any{
		"User":       user,
		"TargetUser": targetUser,
	})
}

// EditUserForm shows edit profile details page
func (r *AdminController) EditUserForm(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	uIDStr := ctx.Request().Route("user_id")
	uID, err := strconv.Atoi(uIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users")
	}

	var targetUser models.User
	err = facades.Orm().Query().Find(&targetUser, uID)
	if err != nil || targetUser.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users")
	}

	return ctx.Response().View().Make("admin/users_edit.tmpl", map[string]any{
		"User":       user,
		"TargetUser": targetUser,
	})
}

// UpdateUser processes updates to user profile details
func (r *AdminController) UpdateUser(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	uIDStr := ctx.Request().Route("user_id")
	uID, err := strconv.Atoi(uIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users")
	}

	var targetUser models.User
	err = facades.Orm().Query().Find(&targetUser, uID)
	if err != nil || targetUser.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users")
	}

	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")
	role := ctx.Request().Input("role")

	if name != "" && email != "" {
		targetUser.Name = name
		targetUser.Email = email
		targetUser.Role = role

		if password != "" {
			hashedPassword, err := facades.Hash().Make(password)
			if err == nil {
				targetUser.Password = hashedPassword
			}
		}

		_ = facades.Orm().Query().Save(&targetUser)
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/users")
}

// DeleteUser deletes a user from administration
func (r *AdminController) DeleteUser(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	uIDStr := ctx.Request().Route("user_id")
	uID, err := strconv.Atoi(uIDStr)
	if err == nil {
		_, _ = facades.Orm().Query().Where("id = ?", uID).Delete(&models.User{})
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/users")
}

// Cms displays CMS panel to manage guest pages
func (r *AdminController) Cms(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var pages []models.CmsPage
	_ = facades.Orm().Query().Order("`order` asc").Get(&pages)

	return ctx.Response().View().Make("admin/cms.tmpl", map[string]any{
		"User":  user,
		"Pages": pages,
	})
}

// CreateCmsPage creates a guest page meta
func (r *AdminController) CreateCmsPage(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	name := ctx.Request().Input("name")
	pType := ctx.Request().Input("type")
	orderStr := ctx.Request().Input("order")
	order, _ := strconv.Atoi(orderStr)

	if name != "" && pType != "" {
		newPage := models.CmsPage{
			Name:  name,
			Type:  pType,
			Order: order,
		}
		_ = facades.Orm().Query().Create(&newPage)
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
}

// UpdateCmsPage updates a guest page properties
func (r *AdminController) UpdateCmsPage(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	pIDStr := ctx.Request().Route("page_id")
	pID, err := strconv.Atoi(pIDStr)
	if err == nil {
		var page models.CmsPage
		_ = facades.Orm().Query().Find(&page, pID)
		if page.ID != 0 {
			page.Name = ctx.Request().Input("name")
			page.Type = ctx.Request().Input("type")
			orderStr := ctx.Request().Input("order")
			order, _ := strconv.Atoi(orderStr)
			page.Order = order
			_ = facades.Orm().Query().Save(&page)
		}
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
}

// DeleteCmsPage deletes page and its contents
func (r *AdminController) DeleteCmsPage(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	pIDStr := ctx.Request().Route("page_id")
	pID, err := strconv.Atoi(pIDStr)
	if err == nil {
		_, _ = facades.Orm().Query().Where("id = ?", pID).Delete(&models.CmsPage{})
		_, _ = facades.Orm().Query().Where("page_id = ?", pID).Delete(&models.GuestContent{})
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
}

// CmsPageDetail shows specific contents inside page grouped by section
func (r *AdminController) CmsPageDetail(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	pIDStr := ctx.Request().Route("page_id")
	pID, err := strconv.Atoi(pIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
	}

	var page models.CmsPage
	_ = facades.Orm().Query().Find(&page, pID)
	if page.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
	}

	var contents []models.GuestContent
	_ = facades.Orm().Query().Where("page_id = ?", page.ID).Get(&contents)

	var products []models.Product
	if page.Type == "products" {
		_ = facades.Orm().Query().Get(&products)
	}

	// Group contents by section to avoid one massive single card in UI
	groupedContents := make(map[string][]models.GuestContent)
	var sectionsOrder []string
	for _, c := range contents {
		sec := c.Section
		if sec == "" {
			sec = "Umum / Lainnya"
		}
		if _, exists := groupedContents[sec]; !exists {
			sectionsOrder = append(sectionsOrder, sec)
		}
		groupedContents[sec] = append(groupedContents[sec], c)
	}

	return ctx.Response().View().Make("admin/cms_detail.tmpl", map[string]any{
		"User":            user,
		"Page":            page,
		"GroupedContents": groupedContents,
		"SectionsOrder":   sectionsOrder,
		"Products":        products,
	})
}

// UpdateProductPrice processes updates to product prices from CMS detail page
func (r *AdminController) UpdateProductPrice(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	prodIDStr := ctx.Request().Route("product_id")
	prodID, err := strconv.Atoi(prodIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
	}

	origPriceStr := ctx.Request().Input("original_price")
	discountStr := ctx.Request().Input("discount")

	origPrice, err1 := strconv.ParseFloat(origPriceStr, 64)
	discount, err2 := strconv.ParseFloat(discountStr, 64)

	if err1 == nil && err2 == nil {
		var product models.Product
		_ = facades.Orm().Query().Find(&product, prodID)
		if product.ID != 0 {
			product.OriginalPrice = origPrice
			product.Discount = discount
			// Calculate discounted price automatically
			product.Price = origPrice * (1 - discount/100)
			_ = facades.Orm().Query().Save(&product)
		}
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/cms/pages/3")
}

// UpdateCmsPageContents updates key values inside page
func (r *AdminController) UpdateCmsPageContents(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	pIDStr := ctx.Request().Route("page_id")
	pID, err := strconv.Atoi(pIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
	}

	var contents []models.GuestContent
	_ = facades.Orm().Query().Where("page_id = ?", pID).Get(&contents)

	for _, c := range contents {
		idKey := c.Key + "_id"
		enKey := c.Key + "_en"
		
		valId := ctx.Request().Input(idKey)
		valEn := ctx.Request().Input(enKey)
		
		hasUpdate := false
		if valId != "" {
			c.ValueId = valId
			hasUpdate = true
		}
		if valEn != "" {
			c.ValueEn = valEn
			hasUpdate = true
		}
		
		if hasUpdate {
			_ = facades.Orm().Query().Save(&c)
		}
	}

	// Handle deferred/batch deleted keys submitted by this form
	deletedKeysStr := ctx.Request().Input("deleted_keys")
	if deletedKeysStr != "" {
		ids := strings.Split(deletedKeysStr, ",")
		for _, idStr := range ids {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				_, _ = facades.Orm().Query().Where("id = ?", id).Delete(&models.GuestContent{})
			}
		}
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/cms/pages/"+pIDStr)
}

// CreateCmsPageContent creates a new translation key for the page
func (r *AdminController) CreateCmsPageContent(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	pIDStr := ctx.Request().Route("page_id")
	pID, err := strconv.Atoi(pIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
	}

	key := ctx.Request().Input("key")
	valId := ctx.Request().Input("value_id")
	valEn := ctx.Request().Input("value_en")
	section := ctx.Request().Input("section")
	style := ctx.Request().Input("style")

	if key != "" {
		newContent := models.GuestContent{
			PageID:  uint(pID),
			Key:     key,
			ValueId: valId,
			ValueEn: valEn,
			Section: section,
			Style:   style,
		}
		_ = facades.Orm().Query().Create(&newContent)
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/cms/pages/"+pIDStr)
}

// DeleteCmsPageContent deletes a specific content translation key
func (r *AdminController) DeleteCmsPageContent(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	contentIDStr := ctx.Request().Route("content_id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err == nil {
		var content models.GuestContent
		_ = facades.Orm().Query().Find(&content, contentID)
		if content.ID != 0 {
			pageIDStr := strconv.Itoa(int(content.PageID))
			_, _ = facades.Orm().Query().Where("id = ?", content.ID).Delete(&models.GuestContent{})
			return ctx.Response().Redirect(http.StatusFound, "/admin/cms/pages/"+pageIDStr)
		}
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/cms")
}

// Roles displays the custom database schemas and roles listing page
func (r *AdminController) Roles(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var roles []models.Role
	_ = facades.Orm().Query().Get(&roles)

	return ctx.Response().View().Make("admin/roles.tmpl", map[string]any{
		"User":  user,
		"Roles": roles,
	})
}

// CreateRole adds a new system role and custom db schema definition
func (r *AdminController) CreateRole(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	name := ctx.Request().Input("name")
	label := ctx.Request().Input("label")
	tableName := ctx.Request().Input("table_name")
	relationName := ctx.Request().Input("relation_name")
	isSystemStr := ctx.Request().Input("is_system", "0")
	isRegisterableStr := ctx.Request().Input("is_registerable", "0")
	badgeColor := ctx.Request().Input("badge_color")
	description := ctx.Request().Input("description")
	dashboardRoute := ctx.Request().Input("dashboard_route")
	dashboardView := ctx.Request().Input("dashboard_view")

	isSystem := isSystemStr == "1" || isSystemStr == "true"
	isRegisterable := isRegisterableStr == "1" || isRegisterableStr == "true"

	if name != "" {
		newRole := models.Role{
			Name:           name,
			Label:          label,
			TableName:      tableName,
			RelationName:   relationName,
			IsSystem:       isSystem,
			IsRegisterable: isRegisterable,
			BadgeColor:     badgeColor,
			Description:    description,
			DashboardRoute: dashboardRoute,
			DashboardView:  dashboardView,
		}
		_ = facades.Orm().Query().Create(&newRole)
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/roles")
}

// UpdateRole modifies an existing system role schema definition
func (r *AdminController) UpdateRole(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	roleIDStr := ctx.Request().Route("role_id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/roles")
	}

	var role models.Role
	err = facades.Orm().Query().Find(&role, roleID)
	if err == nil {
		isSystemStr := ctx.Request().Input("is_system", "0")
		isRegisterableStr := ctx.Request().Input("is_registerable", "0")

		role.Name = ctx.Request().Input("name")
		role.Label = ctx.Request().Input("label")
		role.TableName = ctx.Request().Input("table_name")
		role.RelationName = ctx.Request().Input("relation_name")
		role.IsSystem = isSystemStr == "1" || isSystemStr == "true"
		role.IsRegisterable = isRegisterableStr == "1" || isRegisterableStr == "true"
		role.BadgeColor = ctx.Request().Input("badge_color")
		role.Description = ctx.Request().Input("description")
		role.DashboardRoute = ctx.Request().Input("dashboard_route")
		role.DashboardView = ctx.Request().Input("dashboard_view")
		_ = facades.Orm().Query().Save(&role)
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/roles")
}

// DeleteRole removes a custom role definition from the platform
func (r *AdminController) DeleteRole(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	roleIDStr := ctx.Request().Route("role_id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err == nil {
		var role models.Role
		err = facades.Orm().Query().Find(&role, roleID)
		// Protect critical seed roles
		if err == nil && role.Name != "admin" && role.Name != "client" {
			_, _ = facades.Orm().Query().Delete(&role)
		}
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/roles")
}

// GetDbTables returns list of physical tables in the MySQL database
func (r *AdminController) GetDbTables(ctx http.Context) http.Response {
	var tables []string
	err := facades.Orm().Query().Raw("SHOW TABLES").Scan(&tables)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return ctx.Response().Json(http.StatusOK, tables)
}

// GetDbColumns returns list of columns/fields in a given table
func (r *AdminController) GetDbColumns(ctx http.Context) http.Response {
	tableName := ctx.Request().Route("table")
	if tableName == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{"error": "table parameter is required"})
	}

	type ColumnSchema struct {
		Name       string `gorm:"column:name" json:"name"`
		Type       string `gorm:"column:type" json:"type"`
		FullType   string `gorm:"column:full_type" json:"full_type"`
		IsNullable string `gorm:"column:is_nullable" json:"is_nullable"`
		ColKey     string `gorm:"column:col_key" json:"col_key"`
		Extra      string `gorm:"column:extra" json:"extra"`
		RefTable   string `gorm:"column:ref_table" json:"ref_table"`
		RefColumn  string `gorm:"column:ref_column" json:"ref_column"`
	}
	var cols []ColumnSchema
	query := `
		SELECT 
			c.COLUMN_NAME as name,
			c.DATA_TYPE as type,
			c.COLUMN_TYPE as full_type,
			c.IS_NULLABLE as is_nullable,
			c.COLUMN_KEY as col_key,
			c.EXTRA as extra,
			k.REFERENCED_TABLE_NAME as ref_table,
			k.REFERENCED_COLUMN_NAME as ref_column
		FROM information_schema.COLUMNS c
		LEFT JOIN information_schema.KEY_COLUMN_USAGE k 
			ON c.TABLE_SCHEMA = k.TABLE_SCHEMA 
			AND c.TABLE_NAME = k.TABLE_NAME 
			AND c.COLUMN_NAME = k.COLUMN_NAME
			AND k.REFERENCED_TABLE_NAME IS NOT NULL
		WHERE c.TABLE_SCHEMA = DATABASE() 
			AND c.TABLE_NAME = ?
		ORDER BY c.ORDINAL_POSITION
	`
	err := facades.Orm().Query().Raw(query, tableName).Scan(&cols)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return ctx.Response().Json(http.StatusOK, cols)
}


