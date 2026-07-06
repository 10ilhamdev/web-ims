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

	// Clear existing products to ensure correct IDs and packages
	_, _ = facades.Orm().Query().Model(&models.Product{}).Where("id > 0").ForceDelete()
	var dummy []map[string]any
	_ = facades.Orm().Query().Raw("ALTER TABLE products AUTO_INCREMENT = 1").Scan(&dummy)

	// Seed Products
	products := []models.Product{
		{
			Name:        "Paket Web Standar",
			Description: "Sempurna untuk profil bisnis kecil, portofolio, dan landing page.",
			Price:       1500000,
			Features:    `["Desain Halaman Tunggal Kustom (Landing Page / One-Page Website)","Tampilan Responsif & Estetika Premium","Integrasi Formulir Kontak","Optimasi SEO Dasar","Dukungan Gratis 1 Bulan"]`,
			Image:       "globe",
		},
		{
			Name:        "Paket Bisnis Premium",
			Description: "Website profesional dengan CMS admin panel khusus, blog, dan fungsionalitas dinamis.",
			Price:       4500000,
			Features:    `["5 Halaman Utama Kustom (Beranda, Profil, Layanan, Kontak, Blog)","Postingan/Produk Tidak Terbatas (Melalui CMS Panel)","CMS Admin Panel untuk Mengelola Konten","Integrasi Google Maps & Google Analytics","Pembagian Peran Pengguna (Admin/Staff/Customer)","Dukungan Premium 3 Bulan","Optimasi SEO & Kecepatan Website"]`,
			Image:       "building",
		},
		{
			Name:        "Cetak Biru Portal & Solusi Enterprise Kustom",
			Description: "Platform kustom berskala besar, integrasi API, orkestrasi cloud, dan solusi enterprise yang disesuaikan.",
			Price:       0,
			Features:    `["Arsitektur Kustom Modern (React/Vue + Go/Node/Laravel)","Pengembangan & Integrasi API Kustom","Penyebaran Cloud Mandiri (AWS/GCP/Cloudflare)","Skalabilitas Tinggi & Keamanan Kriptografi","Dukungan SLA Premium 6 Bulan","Pembagian Peran Pengguna Kustom (Admin/Staff/Client)","Jumlah Halaman / Produk Dinamis Tanpa Batas"]`,
			Image:       "cpu",
		},
	}

	for _, p := range products {
		if err := facades.Orm().Query().Create(&p); err != nil {
			return err
		}
	}

	// Seed CMS Pages
	_, _ = facades.Orm().Query().Model(&models.CmsPage{}).Where("id > 0").ForceDelete()
	_ = facades.Orm().Query().Raw("ALTER TABLE cms_pages AUTO_INCREMENT = 1").Scan(&dummy)

	pages := []models.CmsPage{
		{Name: "Beranda", Type: "home", Order: 1},
		{Name: "Layanan", Type: "services", Order: 2},
		{Name: "Produk", Type: "products", Order: 3},
	}
	for i := range pages {
		_ = facades.Orm().Query().Create(&pages[i])
	}

	// Seed CMS Guest Contents linked to Pages
	_, _ = facades.Orm().Query().Model(&models.GuestContent{}).Where("id > 0").ForceDelete()
	_ = facades.Orm().Query().Raw("ALTER TABLE guest_contents AUTO_INCREMENT = 1").Scan(&dummy)

	contents := []models.GuestContent{
		{
			PageID:  1, // Beranda
			Key:     "hero_slogan",
			ValueId: "Mengubah Ide Menjadi Solusi Digital",
			ValueEn: "Turning Ideas Into Digital Solutions",
		},
		{
			PageID:  1, // Beranda
			Key:     "hero_title_1",
			ValueId: "Tingkatkan Bisnis Anda Dengan",
			ValueEn: "Elevate Your Business With",
		},
		{
			PageID:  1, // Beranda
			Key:     "hero_title_2",
			ValueId: "Pengembangan Website Premium",
			ValueEn: "Premium Web Development",
		},
		{
			PageID:  1, // Beranda
			Key:     "hero_desc",
			ValueId: "IMS (Innovation Massive Solutions) merancang dan membangun aplikasi web berkinerja tinggi yang disesuaikan untuk bisnis Anda. Aman, responsif, dan dibuat menggunakan teknologi modern seperti Laravel, React, Vue, Go, dan Node.js.",
			ValueEn: "IMS (Innovation Massive Solutions) designs and engineers high-performance web applications tailored to your business. Secure, responsive, and crafted using modern technologies like Laravel, React, Vue, Go, and Node.js.",
		},
		{
			PageID:  1, // Beranda
			Key:     "why_partner_title",
			ValueId: "Mengapa Bermitra dengan IMS?",
			ValueEn: "Why Partner With IMS?",
		},
		{
			PageID:  1, // Beranda
			Key:     "why_partner_desc",
			ValueId: "Kami membangun sistem web yang unggul dalam kinerja, keamanan, dan pengalaman pengguna.",
			ValueEn: "We build web systems that excel in performance, security, and user experience.",
		},
		{
			PageID:  2, // Layanan
			Key:     "services_slogan",
			ValueId: "Keahlian Kami",
			ValueEn: "Our Expertise",
		},
		{
			PageID:  2, // Layanan
			Key:     "services_title",
			ValueId: "Layanan & Teknologi Komprehensif",
			ValueEn: "Comprehensive Tech Stack & Services",
		},
		{
			PageID:  2, // Layanan
			Key:     "services_desc",
			ValueId: "Kami memanfaatkan kerangka kerja terbaik dan infrastruktur cloud untuk membuat solusi tangguh bagi pertumbuhan bisnis.",
			ValueEn: "We leverage top-tier frameworks and cloud infrastructures to craft robust solutions for business scale.",
		},
		{
			PageID:  3, // Produk
			Key:     "products_slogan",
			ValueId: "Paket & Harga",
			ValueEn: "Packages & Pricing",
		},
		{
			PageID:  3, // Produk
			Key:     "products_title",
			ValueId: "Pilih Paket Pengembangan Anda",
			ValueEn: "Choose Your Development Package",
		},
		{
			PageID:  3, // Produk
			Key:     "products_desc",
			ValueId: "Harga jujur dan transparan. Pilih paket yang sesuai dengan skala bisnis Anda saat ini, dan kami akan membangunnya dengan sempurna.",
			ValueEn: "Honest, transparent pricing. Pick a package that aligns with your current scale, and we'll build it with perfection.",
		},
	}

	for _, c := range contents {
		_ = facades.Orm().Query().Create(&c)
	}

	return nil
}
