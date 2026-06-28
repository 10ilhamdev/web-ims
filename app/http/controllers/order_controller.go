package controllers

import (
	"ims/app/facades"
	"ims/app/models"
	"strconv"

	"github.com/goravel/framework/contracts/http"
)

type OrderController struct{}

func NewOrderController() *OrderController {
	return &OrderController{}
}

// Index lists orders for logged-in user
func (r *OrderController) Index(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var orders []models.Order
	err := facades.Orm().Query().
		Where("user_id = ?", user.ID).
		Get(&orders)
	if err == nil {
		for i := range orders {
			var product models.Product
			_ = facades.Orm().Query().Find(&product, orders[i].ProductID)
			orders[i].Product = &product
		}
	} else {
		orders = []models.Order{}
	}

	return ctx.Response().View().Make("orders.tmpl", map[string]any{
		"User":   user,
		"Orders": orders,
	})
}

// Checkout displays checkout confirmation
func (r *OrderController) Checkout(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	pIDStr := ctx.Request().Route("product_id")
	pID, err := strconv.Atoi(pIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/products")
	}

	var product models.Product
	err = facades.Orm().Query().Find(&product, pID)
	if err != nil || product.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/products")
	}

	return ctx.Response().View().Make("checkout.tmpl", map[string]any{
		"User":    user,
		"Product": product,
	})
}

// PlaceOrder creates order
func (r *OrderController) PlaceOrder(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	pIDStr := ctx.Request().Route("product_id")
	pID, err := strconv.Atoi(pIDStr)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/products")
	}

	var product models.Product
	err = facades.Orm().Query().Find(&product, pID)
	if err != nil || product.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/products")
	}

	requirements := ctx.Request().Input("requirements")
	if requirements == "" {
		return ctx.Response().View().Make("checkout.tmpl", map[string]any{
			"User":    user,
			"Product": product,
			"Error":   "Brief/Requirements is required to place an order.",
		})
	}

	newOrder := models.Order{
		UserID:       user.ID,
		ProductID:    product.ID,
		Requirements: requirements,
		Price:        product.Price,
		Status:       "pending",
	}

	if err := facades.Orm().Query().Create(&newOrder); err != nil {
		return ctx.Response().View().Make("checkout.tmpl", map[string]any{
			"User":    user,
			"Product": product,
			"Error":   "Failed to place order. Please try again.",
		})
	}

	return ctx.Response().Redirect(http.StatusFound, "/orders")
}
