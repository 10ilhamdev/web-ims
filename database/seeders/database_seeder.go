package seeders

import (
	"ims/app/facades"
	"ims/app/models"
)
	
type DatabaseSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *DatabaseSeeder) Signature() string {
	return "DatabaseSeeder"
}

// Run executes the seeder logic.
func (s *DatabaseSeeder) Run() error {
	// Seed Admin User
	count, err := facades.Orm().Query().Model(&models.User{}).Where("email = ?", "admin@ims.com").Count()
	if err != nil {
		return err
	}
	if count == 0 {
		hashedPassword, err := facades.Hash().Make("password123")
		if err != nil {
			return err
		}
		admin := models.User{
			Name:     "Admin IMS",
			Email:    "admin@ims.com",
			Password: hashedPassword,
			Role:     "admin",
		}
		if err := facades.Orm().Query().Create(&admin); err != nil {
			return err
		}
	}

	// Seed Products
	products := []models.Product{
		{
			Name:        "Startup Landing Page",
			Description: "Perfect for new startups wanting to establish an online presence quickly with premium styling.",
			Price:       1500000,
			Features:    `["1 Page Custom Design","Responsive Layout","Contact Form Integration","Basic SEO Optimization","1 Month Free Support"]`,
			Image:       "globe",
		},
		{
			Name:        "Corporate Portal",
			Description: "Complete professional website for businesses needing services showcase, blogs, and CMS integration.",
			Price:       4500000,
			Features:    `["Up to 5 Pages","CMS Admin Panel","Blog / News Section","Google Maps & Analytics","3 Months Support","SEO & Performance Tuning"]`,
			Image:       "building",
		},
		{
			Name:        "E-Commerce Store",
			Description: "Full-featured online shop with product management, shopping cart, checkout, payment gateway, and client panel.",
			Price:       9500000,
			Features:    `["Unlimited Products","Shopping Cart & Checkout","Payment Gateway Integration","Order Dashboard","6 Months Premium Support","Advanced Security Audit"]`,
			Image:       "shopping-cart",
		},
		{
			Name:        "Enterprise Custom Platform",
			Description: "Highly scalable, custom-built web application tailored to your complex enterprise business processes.",
			Price:       20000000,
			Features:    `["Custom Architecture (React/Vue + Go/Node)","API Development & Integrations","Cloud Deployment (AWS/GCP)","High Availability & Scalability","1 Year SLA Support","Custom User Roles & Access"]`,
			Image:       "cpu",
		},
	}

	for _, p := range products {
		pCount, err := facades.Orm().Query().Model(&models.Product{}).Where("name = ?", p.Name).Count()
		if err != nil {
			return err
		}
		if pCount == 0 {
			if err := facades.Orm().Query().Create(&p); err != nil {
				return err
			}
		}
	}

	return nil
}
