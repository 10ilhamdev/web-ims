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
