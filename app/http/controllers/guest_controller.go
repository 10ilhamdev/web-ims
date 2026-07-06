package controllers

import (
	"encoding/json"
	"ims/app/facades"
	"ims/app/models"

	"github.com/goravel/framework/contracts/http"
)

type GuestController struct{}

func NewGuestController() *GuestController {
	return &GuestController{}
}

// Index serves landing page
func (r *GuestController) Index(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	idJson, enJson := GetCmsTranslationsJson(ctx)
	return ctx.Response().View().Make("home.tmpl", map[string]any{
		"User":                  user,
		"CmsTranslationsIDJson": idJson,
		"CmsTranslationsENJson": enJson,
	})
}

// Services serves services detail page
func (r *GuestController) Services(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	idJson, enJson := GetCmsTranslationsJson(ctx)
	return ctx.Response().View().Make("services.tmpl", map[string]any{
		"User":                  user,
		"CmsTranslationsIDJson": idJson,
		"CmsTranslationsENJson": enJson,
	})
}

// Products serves product packages page
func (r *GuestController) Products(ctx http.Context) http.Response {
	user := GetCurrentUser(ctx)
	idJson, enJson := GetCmsTranslationsJson(ctx)
	
	var products []models.Product
	err := facades.Orm().Query().Get(&products)
	if err != nil {
		products = []models.Product{}
	}

	// Prepare data by decoding features JSON for each product
	type ProductView struct {
		models.Product
		FeaturesList []string
	}
	var viewProducts []ProductView
	for _, p := range products {
		var feats []string
		_ = json.Unmarshal([]byte(p.Features), &feats)
		viewProducts = append(viewProducts, ProductView{
			Product:      p,
			FeaturesList: feats,
		})
	}

	return ctx.Response().View().Make("products.tmpl", map[string]any{
		"User":                  user,
		"Products":              viewProducts,
		"CmsTranslationsIDJson": idJson,
		"CmsTranslationsENJson": enJson,
	})
}
