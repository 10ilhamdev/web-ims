package controllers

import (
	"ims/app/facades"
	"ims/app/models"
	"strconv"

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
		return ctx.Response().Redirect(http.StatusFound, "/admin/users")
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

// CmsPageDetail shows specific contents inside page
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

	return ctx.Response().View().Make("admin/cms_detail.tmpl", map[string]any{
		"User":     user,
		"Page":     page,
		"Contents": contents,
	})
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
		valId := ctx.Request().Input(c.Key + "_id")
		valEn := ctx.Request().Input(c.Key + "_en")
		c.ValueId = valId
		c.ValueEn = valEn
		_ = facades.Orm().Query().Save(&c)
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

	if key != "" {
		newContent := models.GuestContent{
			PageID:  uint(pID),
			Key:     key,
			ValueId: valId,
			ValueEn: valEn,
		}
		_ = facades.Orm().Query().Create(&newContent)
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/cms/pages/"+pIDStr)
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
	tableName := ctx.Request().Input("table_name")
	modelName := ctx.Request().Input("model_name")
	fields := ctx.Request().Input("fields")
	relations := ctx.Request().Input("relations")

	if name != "" {
		newRole := models.Role{
			Name:      name,
			TableName: tableName,
			ModelName: modelName,
			Fields:    fields,
			Relations: relations,
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
		role.Name = ctx.Request().Input("name")
		role.TableName = ctx.Request().Input("table_name")
		role.ModelName = ctx.Request().Input("model_name")
		role.Fields = ctx.Request().Input("fields")
		role.Relations = ctx.Request().Input("relations")
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

	type ColumnInfo struct {
		Field string `gorm:"column:Field"`
	}
	var cols []ColumnInfo
	err := facades.Orm().Query().Raw("SHOW COLUMNS FROM " + tableName).Scan(&cols)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	var columns []string
	for _, col := range cols {
		columns = append(columns, col.Field)
	}

	return ctx.Response().Json(http.StatusOK, columns)
}


