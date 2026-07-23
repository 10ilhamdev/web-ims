package seeders

import (
	"ims/app/facades"
	"ims/app/models"
	"strings"
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

	// Seed System Roles (Admin & Client)
	roleCount, err := facades.Orm().Query().Model(&models.Role{}).Count()
	if err == nil && roleCount == 0 {
		adminRole := models.Role{
			Name:           "admin",
			Label:          "Administrator",
			TableName:      "admins",
			RelationName:   "Admin",
			IsSystem:       true,
			IsRegisterable: false,
			BadgeColor:     "danger",
			Description:    "Platform main administrator with all management privileges.",
			DashboardRoute: "/admin",
			DashboardView:  "admin/dashboard",
		}
		_ = facades.Orm().Query().Create(&adminRole)

		clientRole := models.Role{
			Name:           "client",
			Label:          "Client / Customer",
			TableName:      "customers",
			RelationName:   "Customer",
			IsSystem:       false,
			IsRegisterable: true,
			BadgeColor:     "success",
			Description:    "Company clients and customers who place website package development orders.",
			DashboardRoute: "/dashboard",
			DashboardView:  "dashboard",
		}
		_ = facades.Orm().Query().Create(&clientRole)
	}

	// Clear existing products to ensure correct IDs and packages
	_, _ = facades.Orm().Query().Model(&models.Product{}).Where("id > 0").ForceDelete()
	var dummy []map[string]any
	_ = facades.Orm().Query().Raw("ALTER TABLE products AUTO_INCREMENT = 1").Scan(&dummy)

	// Seed Products
	products := []models.Product{
		{
			Name:          "Paket Web Standar",
			Description:   "Sempurna untuk profil bisnis kecil, portofolio, dan landing page.",
			Price:         500000,
			OriginalPrice: 1000000,
			Discount:      50,
			Features:      `["Desain Halaman Tunggal Kustom (Landing Page / One-Page Website)","Tampilan Responsif & Estetika Premium","Integrasi Formulir Kontak","Optimasi SEO Dasar","Dukungan Gratis 1 Bulan"]`,
			Image:         "globe",
		},
		{
			Name:          "Paket Bisnis Premium",
			Description:   "Website profesional dengan CMS admin panel khusus, blog, dan fungsionalitas dinamis.",
			Price:         4500000,
			OriginalPrice: 8000000,
			Discount:      43.75,
			Features:      `["5 Halaman Utama Kustom (Beranda, Profil, Layanan, Kontak, Blog)","Postingan/Produk Tidak Terbatas (Melalui CMS Panel)","CMS Admin Panel untuk Mengelola Konten","Integrasi Google Maps & Google Analytics","Pembagian Peran Pengguna (Admin/Staff/Customer)","Dukungan Premium 3 Bulan","Optimasi SEO & Kecepatan Website"]`,
			Image:         "building",
		},
		{
			Name:          "Cetak Biru Portal & Solusi Enterprise Kustom",
			Description:   "Platform kustom berskala besar, integrasi API, orkestrasi cloud, dan solusi enterprise yang disesuaikan.",
			Price:         0,
			OriginalPrice: 0,
			Discount:      0,
			Features:      `["Arsitektur Kustom Modern (React/Vue + Go/Node/Laravel)","Pengembangan & Integrasi API Kustom","Penyebaran Cloud Mandiri (AWS/GCP/Cloudflare)","Skalabilitas Tinggi & Keamanan Kriptografi","Dukungan SLA Premium 6 Bulan","Pembagian Peran Pengguna Kustom (Admin/Staff/Client)","Jumlah Halaman / Produk Dinamis Tanpa Batas"]`,
			Image:         "cpu",
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
		{Name: "Navigasi & Umum", Type: "general", Order: 4},
	}
	for i := range pages {
		_ = facades.Orm().Query().Create(&pages[i])
	}

	// Seed CMS Guest Contents linked to Pages
	_, _ = facades.Orm().Query().Model(&models.GuestContent{}).Where("id > 0").ForceDelete()
	_ = facades.Orm().Query().Raw("ALTER TABLE guest_contents AUTO_INCREMENT = 1").Scan(&dummy)

	contents := []models.GuestContent{
		// PageID 1: Beranda
		{
			PageID:  1,
			Key:     "hero_slogan",
			ValueId: "Mengubah Ide Menjadi Solusi Digital",
			ValueEn: "Turning Ideas Into Digital Solutions",
		},
		{
			PageID:  1,
			Key:     "hero_title_1",
			ValueId: "Tingkatkan Bisnis Anda Dengan",
			ValueEn: "Grow Your Business With",
		},
		{
			PageID:  1,
			Key:     "hero_title_2",
			ValueId: "Pengembangan Website Premium",
			ValueEn: "Premium Website Development",
		},
		{
			PageID:  1,
			Key:     "hero_desc",
			ValueId: "IMS (Innovation Massive Solutions) merancang dan membangun aplikasi web berkinerja tinggi yang disesuaikan untuk bisnis Anda. Aman, responsif, dan dibuat menggunakan teknologi modern seperti Laravel, React, Vue, Go, dan Node.js.",
			ValueEn: "IMS (Innovation Massive Solutions) designs and builds high-performance web applications tailored to your business. Secure, responsive, and crafted using modern tech like Laravel, React, Vue, Go, and Node.js.",
		},
		{
			PageID:  1,
			Key:     "btn_view_packages",
			ValueId: "Lihat Paket Kami",
			ValueEn: "View Our Packages",
		},
		{
			PageID:  1,
			Key:     "btn_explore_services",
			ValueId: "Jelajahi Layanan",
			ValueEn: "Explore Services",
		},
		{
			PageID:  1,
			Key:     "stack_banner_text",
			ValueId: "Teknologi yang Digunakan",
			ValueEn: "Technologies We Leverage",
		},
		{
			PageID:  1,
			Key:     "why_partner_title",
			ValueId: "Mengapa Bermitra dengan IMS?",
			ValueEn: "Why Partner With IMS?",
		},
		{
			PageID:  1,
			Key:     "why_partner_desc",
			ValueId: "Kami membangun sistem web yang unggul dalam kinerja, keamanan, dan pengalaman pengguna.",
			ValueEn: "We engineer web systems that excel in speed, protection, and user experience.",
		},
		{
			PageID:  1,
			Key:     "feature_clean_code",
			ValueId: "Kode Bersih Modern",
			ValueEn: "Modern Clean Code",
		},
		{
			PageID:  1,
			Key:     "feature_clean_code_desc",
			ValueId: "Memanfaatkan arsitektur perangkat lunak yang kuat dengan React, Node.js, dan Golang untuk memastikan pemeliharaan basis kode yang mudah.",
			ValueEn: "Leveraging robust software architectures with React, Node.js, and Golang to ensure easy codebase maintenance.",
		},
		{
			PageID:  1,
			Key:     "feature_security",
			ValueId: "Keamanan Skala Enterprise",
			ValueEn: "Enterprise-Grade Security",
		},
		{
			PageID:  1,
			Key:     "feature_security_desc",
			ValueId: "Kami menerapkan enkripsi tingkat tinggi, pemeriksaan keamanan digital, dan mekanisme otentikasi yang aman di setiap proyek.",
			ValueEn: "We deploy state-of-the-art encryption, digital signature logic, and secure authentication schemas on every single project.",
		},
		{
			PageID:  1,
			Key:     "feature_performance",
			ValueId: "Performa Tinggi",
			ValueEn: "Blazing Fast Performance",
		},
		{
			PageID:  1,
			Key:     "feature_performance_desc",
			ValueId: "Pemuatan halaman yang sangat cepat dengan optimasi sisi server, cache konten, dan sistem desain yang responsif.",
			ValueEn: "Hyper-optimized page speed through advanced caching, light assets, and responsive custom layout architectures.",
		},
		{
			PageID:  1,
			Key:     "cta_ready_title",
			ValueId: "Siap Meluncurkan Usaha Digital Anda?",
			ValueEn: "Ready to Unleash Your Digital Product?",
		},
		{
			PageID:  1,
			Key:     "cta_ready_desc",
			ValueId: "Pilih dari paket web yang sudah kami siapkan atau ajukan cetak biru portal kustom yang dirancang sesuai kebutuhan Anda.",
			ValueEn: "Select a prepared configuration or request a customized blueprint matching your requirements.",
		},
		{
			PageID:  1,
			Key:     "cta_btn_get_started",
			ValueId: "Mulai Sekarang",
			ValueEn: "Get Started Now",
		},

		// PageID 2: Layanan
		{
			PageID:  2,
			Key:     "services_slogan",
			ValueId: "Keahlian Kami",
			ValueEn: "Our Expertise",
		},
		{
			PageID:  2,
			Key:     "services_title",
			ValueId: "Layanan & Teknologi Komprehensif",
			ValueEn: "Comprehensive Services & Tech Stack",
		},
		{
			PageID:  2,
			Key:     "services_desc",
			ValueId: "Kami memanfaatkan kerangka kerja terbaik dan infrastruktur cloud untuk membuat solusi tangguh bagi pertumbuhan bisnis.",
			ValueEn: "We harness industry-best frameworks and cloud infrastructure to create resilient solutions for business growth.",
		},
		{
			PageID:  2,
			Key:     "service_laravel_title",
			ValueId: "Pengembangan Laravel",
			ValueEn: "Laravel Engineering",
		},
		{
			PageID:  2,
			Key:     "service_laravel_desc",
			ValueId: "Membangun aplikasi perusahaan berbasis database dan modular. Sangat cocok untuk portal, panel administrasi, dan API backend yang tangguh.",
			ValueEn: "Building modular, database-driven business applications. Perfect for portals, administrative panels, and robust backend APIs.",
		},
		{
			PageID:  2,
			Key:     "service_react_title",
			ValueId: "React & SPA Platforms",
			ValueEn: "React & SPA Platforms",
		},
		{
			PageID:  2,
			Key:     "service_react_desc",
			ValueId: "Dasbor klien berkinerja tinggi, antarmuka dinamis, dan aplikasi satu halaman (SPA) menggunakan Next.js dan manajemen status kustom.",
			ValueEn: "High-performance client dashboards, dynamic interfaces, and single-page apps (SPA) utilizing Next.js and custom state management.",
		},
		{
			PageID:  2,
			Key:     "service_vue_title",
			ValueId: "Aplikasi Vue.js & Nuxt.js",
			ValueEn: "Vue.js & Nuxt.js Apps",
		},
		{
			PageID:  2,
			Key:     "service_vue_desc",
			ValueId: "Pengalaman pengguna yang interaktif dengan frontend reaktif. Pembuatan prototipe cepat dan layanan integrasi web yang mulus.",
			ValueEn: "Interactive user experiences with reactive frontends. Rapid prototyping and smooth web integration services.",
		},
		{
			PageID:  2,
			Key:     "service_node_title",
			ValueId: "Arsitektur Server Node.js",
			ValueEn: "Node.js Server Architecture",
		},
		{
			PageID:  2,
			Key:     "service_node_desc",
			ValueId: "API RESTful & GraphQL yang skalabel, server komunikasi real-time (WebSockets), dan konfigurasi layanan mikro.",
			ValueEn: "Scalable RESTful & GraphQL APIs, real-time communication servers (WebSockets), and microservices configuration.",
		},
		{
			PageID:  2,
			Key:     "service_python_title",
			ValueId: "Otomatisasi Python & AI",
			ValueEn: "Python Automation & AI",
		},
		{
			PageID:  2,
			Key:     "service_python_desc",
			ValueId: "Saluran pemrosesan data, integrasi pembelajaran mesin, analitik kecerdasan bisnis, dan sistem otomatisasi latar belakang.",
			ValueEn: "Data processing pipelines, machine learning integration, business intelligence analytics, and background automation systems.",
		},
		{
			PageID:  2,
			Key:     "service_wordpress_title",
			ValueId: "Portal Premium WordPress",
			ValueEn: "WordPress Premium Portals",
		},
		{
			PageID:  2,
			Key:     "service_wordpress_desc",
			ValueId: "Custom landing pages, premium corporate websites, blogs, and content systems with optimal speed and security tuning.",
			ValueEn: "Custom landing pages, premium corporate websites, blogs, and content systems with optimal speed and security tuning.",
		},
		{
			PageID:  2,
			Key:     "service_cloud_title",
			ValueId: "Infrastruktur Cloud",
			ValueEn: "Cloud Infrastructure",
		},
		{
			PageID:  2,
			Key:     "service_cloud_desc",
			ValueId: "Orkestrasi Docker, penyebaran AWS/GCP, alur integrasi dan pengiriman berkelanjutan (CI/CD), dan pemantauan server.",
			ValueEn: "Docker orchestration, AWS/GCP deployment, continuous integration and delivery (CI/CD) pipelines, and server monitoring.",
		},
		{
			PageID:  2,
			Key:     "service_security_title",
			ValueId: "Keamanan & Kriptografi",
			ValueEn: "Security & Cryptography",
		},
		{
			PageID:  2,
			Key:     "service_security_desc",
			ValueId: "Enkripsi ujung-ke-ujung, logika tanda tangan digital, audit basis data yang aman, dan penilaian kepatuhan keamanan.",
			ValueEn: "End-to-end encryption, digital signature logic, secure database audits, and security compliance assessments.",
		},
		{
			PageID:  2,
			Key:     "service_golang_title",
			ValueId: "Pengembangan Backend Golang",
			ValueEn: "Golang Backend Development",
		},
		{
			PageID:  2,
			Key:     "service_golang_desc",
			ValueId: "Membangun microservices berkinerja tinggi, API gRPC/REST, dan pemrosesan paralel cepat dengan keamanan memori optimal.",
			ValueEn: "Building high-performance microservices, fast gRPC/REST APIs, and concurrent processing pipelines with memory-safety.",
		},
		{
			PageID:  2,
			Key:     "service_integration_title",
			ValueId: "Integrasi Layanan Google & CDN",
			ValueEn: "Google Integrations & CDN",
		},
		{
			PageID:  2,
			Key:     "service_integration_desc",
			ValueId: "Penyimpanan cloud Google Drive, pelacakan trafik dengan Google Analytics, dan optimasi kecepatan dengan CDN Cloudflare.",
			ValueEn: "Google Drive cloud storage integration, traffic tracking with Google Analytics, and edge performance optimization via Cloudflare CDN.",
		},
		{
			PageID:  2,
			Key:     "service_charting_title",
			ValueId: "Visualisasi Data & Grafik",
			ValueEn: "Data Visualization & Charts",
		},
		{
			PageID:  2,
			Key:     "service_charting_desc",
			ValueId: "Pembuatan dashboard analitik interaktif menggunakan CDN pustaka grafik (Chart.js, ApexCharts) untuk pie chart, bar chart, dan line chart.",
			ValueEn: "Creating interactive analytic dashboards using graphing libraries (Chart.js, ApexCharts) for pie charts, bar charts, and line charts.",
		},
		{
			PageID:  2,
			Key:     "service_vr_title",
			ValueId: "Ruangan Virtual & 3D WebGL",
			ValueEn: "Virtual Rooms & 3D WebGL",
		},
		{
			PageID:  2,
			Key:     "service_vr_desc",
			ValueId: "Membangun tur virtual 3D, showroom interaktif, dan ruang virtual imersif berbasis Three.js, A-Frame, dan WebGL.",
			ValueEn: "Building immersive 3D virtual tours, interactive showrooms, and visual spaces powered by Three.js, A-Frame, and WebGL.",
		},
		{
			PageID:  2,
			Key:     "service_flipbook_title",
			ValueId: "Flipbook & Media Interaktif",
			ValueEn: "Flipbook & Interactive Media",
		},
		{
			PageID:  2,
			Key:     "service_flipbook_desc",
			ValueId: "Konversi PDF ke flipbook interaktif 3D (Turn.js, dFlip) dengan efek balik halaman realistis untuk katalog, brosur, dan majalah digital.",
			ValueEn: "Converting PDFs into 3D interactive flipbooks (Turn.js, dFlip) with realistic page-flip effects for catalogs, brochures, and digital magazines.",
		},
		{
			PageID:  2,
			Key:     "service_realtime_title",
			ValueId: "Komunikasi Real-Time & WebRTC",
			ValueEn: "Real-Time Comms & WebRTC",
		},
		{
			PageID:  2,
			Key:     "service_realtime_desc",
			ValueId: "Mengintegrasikan WebRTC dan WebSockets untuk konferensi video peer-to-peer, chat real-time, kolaborasi multi-pengguna, dan live streaming.",
			ValueEn: "Integrating WebRTC and WebSockets for peer-to-peer video conferencing, live chat systems, multi-user collaboration, and broadcasting.",
		},
		{
			PageID:  2,
			Key:     "service_pwa_title",
			ValueId: "Progressive Web App (PWA)",
			ValueEn: "Progressive Web Apps (PWA)",
		},
		{
			PageID:  2,
			Key:     "service_pwa_desc",
			ValueId: "Menciptakan aplikasi web yang dapat diinstal langsung di perangkat user, berjalan secara offline (Service Workers), dan mendukung push notifications.",
			ValueEn: "Designing web applications installable directly on devices, running offline (Service Workers), and supporting push alerts.",
		},
		{
			PageID:  2,
			Key:     "service_wasm_title",
			ValueId: "High-Performance WebAssembly (Wasm)",
			ValueEn: "High-Performance WebAssembly (Wasm)",
		},
		{
			PageID:  2,
			Key:     "service_wasm_desc",
			ValueId: "Kompilasi kode C++/Rust ke Wasm untuk memproses data berat langsung di browser, seperti editor video/gambar dan game web.",
			ValueEn: "Compiling C++/Rust to Wasm for heavy CPU processing inside the browser, like video/image editing and web-based games.",
		},
		{
			PageID:  2,
			Key:     "services_custom_title",
			ValueId: "Punya Proyek Kustom?",
			ValueEn: "Need a Custom Project?",
		},
		{
			PageID:  2,
			Key:     "services_custom_desc",
			ValueId: "Jika kebutuhan Anda di luar paket konfigurasi kami, kami dapat menyusun cetak biru kustom yang dirancang sesuai arsitektur target Anda.",
			ValueEn: "If your requirements are beyond our pre-packaged configurations, we can draft a custom blueprint tailored to your target architecture.",
		},
		{
			PageID:  2,
			Key:     "services_custom_btn",
			ValueId: "Hubungi Arsitek Kami",
			ValueEn: "Contact Our Architect",
		},

		// PageID 3: Produk
		{
			PageID:  3,
			Key:     "products_slogan",
			ValueId: "Paket & Harga",
			ValueEn: "Packages & Pricing",
		},
		{
			PageID:  3,
			Key:     "products_title",
			ValueId: "Pilih Paket Pengembangan Anda",
			ValueEn: "Choose Your Development Package",
		},
		{
			PageID:  3,
			Key:     "products_desc",
			ValueId: "Harga jujur dan transparan. Pilih paket yang sesuai dengan skala bisnis Anda saat ini, dan kami akan membangunnya dengan sempurna.",
			ValueEn: "Honest, transparent pricing. Pick a package that aligns with your current scale, and we'll build it with perfection.",
		},
		{
			PageID:  3,
			Key:     "products_whats_included",
			ValueId: "Yang Termasuk",
			ValueEn: "What's Included",
		},
		{
			PageID:  3,
			Key:     "product_btn_1",
			ValueId: "Pesan Sekarang",
			ValueEn: "Order Now",
		},
		{
			PageID:  3,
			Key:     "product_btn_2",
			ValueId: "Pesan Sekarang",
			ValueEn: "Order Now",
		},
		{
			PageID:  3,
			Key:     "product_btn_3",
			ValueId: "Hubungi Kami",
			ValueEn: "Get In Touch",
		},
		{
			PageID:  3,
			Key:     "product_price_negotiable",
			ValueId: "Hubungi Kami",
			ValueEn: "Contact Us",
		},
		{
			PageID:  3,
			Key:     "product_name_1",
			ValueId: "Paket Web Standar",
			ValueEn: "Standard Web Package",
		},
		{
			PageID:  3,
			Key:     "product_name_2",
			ValueId: "Paket Bisnis Premium",
			ValueEn: "Premium Business Package",
		},
		{
			PageID:  3,
			Key:     "product_name_3",
			ValueId: "Cetak Biru Portal & Solusi Enterprise Kustom",
			ValueEn: "Corporate Portal Blueprint & Custom Enterprise Solution",
		},
		{
			PageID:  3,
			Key:     "product_desc_1",
			ValueId: "Sempurna untuk profil bisnis kecil, portofolio, dan landing page.",
			ValueEn: "Perfect for small business profiles, portfolios, and landing pages.",
		},
		{
			PageID:  3,
			Key:     "product_desc_2",
			ValueId: "Website profesional dengan CMS admin panel khusus, blog, dan fungsionalitas dinamis.",
			ValueEn: "Professional website with dedicated CMS admin panel, blog, and dynamic features.",
		},
		{
			PageID:  3,
			Key:     "product_desc_3",
			ValueId: "Platform kustom berskala besar, integrasi API, orkestrasi cloud, dan solusi enterprise yang disesuaikan.",
			ValueEn: "Large-scale custom platform, API integrations, cloud orchestration, and tailored enterprise solutions.",
		},
		{
			PageID:  3,
			Key:     "feature_standard_1",
			ValueId: "Desain Halaman Tunggal Kustom (Landing Page / One-Page Website)",
			ValueEn: "Custom Single-Page Design (Landing Page / One-Page Website)",
		},
		{
			PageID:  3,
			Key:     "feature_standard_2",
			ValueId: "Tampilan Responsif & Estetika Premium",
			ValueEn: "Responsive Layout & Premium Aesthetics",
		},
		{
			PageID:  3,
			Key:     "feature_standard_3",
			ValueId: "Integrasi Formulir Kontak",
			ValueEn: "Contact Form Integration",
		},
		{
			PageID:  3,
			Key:     "feature_standard_4",
			ValueId: "Optimasi SEO Dasar",
			ValueEn: "Basic SEO Optimization",
		},
		{
			PageID:  3,
			Key:     "feature_standard_5",
			ValueId: "Dukungan Gratis 1 Bulan",
			ValueEn: "1-Month Free Support",
		},
		{
			PageID:  3,
			Key:     "feature_premium_1",
			ValueId: "5 Halaman Utama Kustom (Beranda, Profil, Layanan, Kontak, Blog)",
			ValueEn: "5 Custom Main Pages (Home, Profile, Services, Contact, Blog)",
		},
		{
			PageID:  3,
			Key:     "feature_premium_2",
			ValueId: "Postingan/Produk Tidak Terbatas (Melalui CMS Panel)",
			ValueEn: "Unlimited Posts/Products (Via CMS Panel)",
		},
		{
			PageID:  3,
			Key:     "feature_premium_3",
			ValueId: "CMS Admin Panel untuk Mengelola Konten",
			ValueEn: "CMS Admin Panel for Content Management",
		},
		{
			PageID:  3,
			Key:     "feature_premium_4",
			ValueId: "Integrasi Google Maps & Google Analytics",
			ValueEn: "Google Maps & Google Analytics Integration",
		},
		{
			PageID:  3,
			Key:     "feature_premium_5",
			ValueId: "Pembagian Peran Pengguna (Admin/Staff/Customer)",
			ValueEn: "User Role Division (Admin/Staff/Customer)",
		},
		{
			PageID:  3,
			Key:     "feature_premium_6",
			ValueId: "Dukungan Premium 3 Bulan",
			ValueEn: "3-Month Premium Support",
		},
		{
			PageID:  3,
			Key:     "feature_premium_7",
			ValueId: "Optimasi SEO & Kecepatan Website",
			ValueEn: "SEO & Website Speed Optimization",
		},
		{
			PageID:  3,
			Key:     "feature_custom_1",
			ValueId: "Arsitektur Kustom Modern (React/Vue + Go/Node/Laravel)",
			ValueEn: "Modern Custom Architecture (React/Vue + Go/Node/Laravel)",
		},
		{
			PageID:  3,
			Key:     "feature_custom_2",
			ValueId: "Pengembangan & Integrasi API Kustom",
			ValueEn: "Custom API Development & Integration",
		},
		{
			PageID:  3,
			Key:     "feature_custom_3",
			ValueId: "Penyebaran Cloud Mandiri (AWS/GCP/Cloudflare)",
			ValueEn: "Self-Managed Cloud Deployment (AWS/GCP/Cloudflare)",
		},
		{
			PageID:  3,
			Key:     "feature_custom_4",
			ValueId: "Skalabilitas Tinggi & Keamanan Kriptografi",
			ValueEn: "High Scalability & Cryptographic Security",
		},
		{
			PageID:  3,
			Key:     "feature_custom_5",
			ValueId: "Dukungan SLA Premium 6 Bulan",
			ValueEn: "6-Month Premium SLA Support",
		},
		{
			PageID:  3,
			Key:     "feature_custom_6",
			ValueId: "Pembagian Peran Pengguna Kustom (Admin/Staff/Client)",
			ValueEn: "Custom User Role Division (Admin/Staff/Client)",
		},
		{
			PageID:  3,
			Key:     "feature_custom_7",
			ValueId: "Jumlah Halaman / Produk Dinamis Tanpa Batas",
			ValueEn: "Unlimited Dynamic Pages / Products",
		},

		// PageID 4: Navigasi & Umum
		{
			PageID:  4,
			Key:     "nav_home",
			ValueId: "Beranda",
			ValueEn: "Home",
		},
		{
			PageID:  4,
			Key:     "nav_services",
			ValueId: "Layanan",
			ValueEn: "Services",
		},
		{
			PageID:  4,
			Key:     "nav_products",
			ValueId: "Produk",
			ValueEn: "Products",
		},
		{
			PageID:  4,
			Key:     "nav_login",
			ValueId: "Masuk",
			ValueEn: "Login",
		},
		{
			PageID:  4,
			Key:     "nav_register",
			ValueId: "Daftar",
			ValueEn: "Register",
		},
		{
			PageID:  4,
			Key:     "nav_admin",
			ValueId: "Panel Admin",
			ValueEn: "Admin Panel",
		},
		{
			PageID:  4,
			Key:     "nav_my_orders",
			ValueId: "Pesanan Saya",
			ValueEn: "My Orders",
		},
		{
			PageID:  4,
			Key:     "nav_logout",
			ValueId: "Keluar",
			ValueEn: "Logout",
		},
		{
			PageID:  4,
			Key:     "hi_user",
			ValueId: "Halo, ",
			ValueEn: "Hi, ",
		},
		{
			PageID:  4,
			Key:     "footer_desc",
			ValueId: "Innovation Massive Solutions (IMS) menyediakan layanan pengembangan website profesional dengan beragam teknologi modern. Kami menghadirkan solusi inovatif, aman, skalabel, dan ramah pengguna untuk mendukung transformasi bisnis klien.",
			ValueEn: "Innovation Massive Solutions (IMS) provides professional website development using cutting-edge technologies. We deliver secure, scalable, and responsive digital solutions to support your business transformation.",
		},
		{
			PageID:  4,
			Key:     "footer_nav_title",
			ValueId: "Navigasi",
			ValueEn: "Navigation",
		},
		{
			PageID:  4,
			Key:     "footer_contact_title",
			ValueId: "Hubungi Kami",
			ValueEn: "Contact Us",
		},
		{
			PageID:  4,
			Key:     "footer_copy",
			ValueId: "© 2026 Innovation Massive Solutions (IMS). Hak Cipta Dilindungi Undang-Undang.",
			ValueEn: "© 2026 Innovation Massive Solutions (IMS). All rights reserved.",
		},
		{
			PageID:  4,
			Key:     "login_title",
			ValueId: "Masuk",
			ValueEn: "Sign In",
		},
		{
			PageID:  4,
			Key:     "login_desc",
			ValueId: "Akses akun & dashboard IMS Anda",
			ValueEn: "Access your IMS account & dashboard",
		},
		{
			PageID:  4,
			Key:     "email_label",
			ValueId: "Alamat Email",
			ValueEn: "Email Address",
		},
		{
			PageID:  4,
			Key:     "password_label",
			ValueId: "Kata Sandi",
			ValueEn: "Password",
		},
		{
			PageID:  4,
			Key:     "btn_login",
			ValueId: "Masuk",
			ValueEn: "Login",
		},
		{
			PageID:  4,
			Key:     "no_account",
			ValueId: "Belum punya akun?",
			ValueEn: "Don't have an account?",
		},
		{
			PageID:  4,
			Key:     "register_here",
			ValueId: "Daftar di sini",
			ValueEn: "Register here",
		},
		{
			PageID:  4,
			Key:     "register_title",
			ValueId: "Daftar Akun",
			ValueEn: "Create Account",
		},
		{
			PageID:  4,
			Key:     "register_desc",
			ValueId: "Mulai perjalanan transformasi digital Anda bersama IMS",
			ValueEn: "Start your digital transformation journey with IMS",
		},
		{
			PageID:  4,
			Key:     "name_label",
			ValueId: "Nama Lengkap",
			ValueEn: "Full Name",
		},
		{
			PageID:  4,
			Key:     "btn_register",
			ValueId: "Daftar",
			ValueEn: "Register",
		},
		{
			PageID:  4,
			Key:     "already_account",
			ValueId: "Sudah punya akun?",
			ValueEn: "Already have an account?",
		},
		{
			PageID:  4,
			Key:     "signin_here",
			ValueId: "Masuk di sini",
			ValueEn: "Sign In here",
		},
		{
			PageID:  4,
			Key:     "chat_ai_assistant",
			ValueId: "Asisten AI",
			ValueEn: "AI Assistant",
		},
		{
			PageID:  4,
			Key:     "chat_with_admin",
			ValueId: "Hubungi Admin",
			ValueEn: "Chat with Admin",
		},
		{
			PageID:  4,
			Key:     "chat_greeting",
			ValueId: "Halo! Saya asisten digital IMS. Pilih topik di bawah atau ketik pertanyaan Anda!",
			ValueEn: "Hello! I am your IMS digital assistant. Select a topic or type your query below!",
		},
		{
			PageID:  4,
			Key:     "chat_faq_packages",
			ValueId: "📦 Paket & Harga",
			ValueEn: "📦 Pricing Packages",
		},
		{
			PageID:  4,
			Key:     "chat_faq_services",
			ValueId: "🛠️ Layanan",
			ValueEn: "🛠️ Services",
		},
		{
			PageID:  4,
			Key:     "chat_faq_contact",
			ValueId: "📞 Kontak",
			ValueEn: "📞 Contact Info",
		},
		{
			PageID:  4,
			Key:     "chat_placeholder_ai",
			ValueId: "Tanya AI...",
			ValueEn: "Ask AI...",
		},
		{
			PageID:  4,
			Key:     "chat_placeholder_admin",
			ValueId: "Pesan ke admin...",
			ValueEn: "Message admin...",
		},
		{
			PageID:  4,
			Key:     "chat_support_desc",
			ValueId: "Ketik nama dan email Anda untuk memulai sesi obrolan langsung dengan tim kami.",
			ValueEn: "Type your name and email to start a direct support session with our team.",
		},
		{
			PageID:  4,
			Key:     "chat_support_name",
			ValueId: "Nama Anda",
			ValueEn: "Your Name",
		},
		{
			PageID:  4,
			Key:     "chat_support_email",
			ValueId: "Email Anda",
			ValueEn: "Your Email",
		},
		{
			PageID:  4,
			Key:     "chat_support_start",
			ValueId: "Mulai Obrolan",
			ValueEn: "Start Live Chat",
		},
	}

	for _, c := range contents {
		c.Section = getSectionForKey(c.PageID, c.Key)
		c.Style = getStyleForKey(c.Key)
		_ = facades.Orm().Query().Create(&c)
	}

	return nil
}

func getStyleForKey(key string) string {
	key = strings.ToLower(key)
	if strings.Contains(key, "btn_") || strings.Contains(key, "_btn") || strings.Contains(key, "button") {
		return "button"
	}
	if strings.Contains(key, "title") || strings.Contains(key, "slogan") || strings.Contains(key, "name") {
		return "title"
	}
	if strings.Contains(key, "desc") || strings.Contains(key, "copy") || strings.Contains(key, "description") {
		return "text"
	}
	if strings.Contains(key, "feature_") || strings.Contains(key, "_feature") || strings.Contains(key, "features") {
		return "feature"
	}
	return "general"
}

func getSectionForKey(pageID uint, key string) string {
	switch pageID {
	case 1: // Beranda
		if strings.HasPrefix(key, "hero_") || strings.HasPrefix(key, "btn_") {
			return "Hero Banner"
		}
		if strings.HasPrefix(key, "stack_") {
			return "Teknologi (Stack)"
		}
		if strings.HasPrefix(key, "why_") || strings.HasPrefix(key, "feature_") {
			return "Kelebihan IMS (Why Us)"
		}
		if strings.HasPrefix(key, "cta_") {
			return "Call To Action (CTA)"
		}
	case 2: // Layanan
		if strings.HasPrefix(key, "services_") {
			return "Header Halaman"
		}
		if strings.HasPrefix(key, "service_laravel_") {
			return "Layanan: Laravel"
		}
		if strings.HasPrefix(key, "service_react_") {
			return "Layanan: React & SPA"
		}
		if strings.HasPrefix(key, "service_vue_") {
			return "Layanan: Vue & Nuxt"
		}
		if strings.HasPrefix(key, "service_node_") {
			return "Layanan: Node.js API"
		}
		if strings.HasPrefix(key, "service_python_") {
			return "Layanan: Python & AI"
		}
		if strings.HasPrefix(key, "service_wordpress_") {
			return "Layanan: WordPress"
		}
		if strings.HasPrefix(key, "service_cloud_") {
			return "Layanan: Cloud & DevOps"
		}
		if strings.HasPrefix(key, "service_security_") {
			return "Layanan: Security"
		}
		if strings.HasPrefix(key, "service_golang_") {
			return "Layanan: Golang Backend"
		}
		if strings.HasPrefix(key, "service_integration_") {
			return "Layanan: Google & CDN"
		}
		if strings.HasPrefix(key, "service_charting_") {
			return "Layanan: Data Viz & Charts"
		}
		if strings.HasPrefix(key, "service_vr_") {
			return "Layanan: VR & 3D WebGL"
		}
		if strings.HasPrefix(key, "service_flipbook_") {
			return "Layanan: Flipbook Interactive"
		}
		if strings.HasPrefix(key, "service_realtime_") {
			return "Layanan: Real-Time & WebRTC"
		}
		if strings.HasPrefix(key, "service_pwa_") {
			return "Layanan: PWA"
		}
		if strings.HasPrefix(key, "service_wasm_") {
			return "Layanan: WebAssembly"
		}
		if strings.HasPrefix(key, "services_custom_") {
			return "Proyek Kustom (Custom Service)"
		}
	case 3: // Produk
		if strings.HasPrefix(key, "products_slogan") || strings.HasPrefix(key, "products_title") || strings.HasPrefix(key, "products_desc") || strings.HasPrefix(key, "products_whats_") {
			return "Header Halaman"
		}
		if strings.HasSuffix(key, "_1") || strings.HasPrefix(key, "feature_standard_") || key == "product_btn_1" {
			return "Paket 1: Paket Web Standar"
		}
		if strings.HasSuffix(key, "_2") || strings.HasPrefix(key, "feature_premium_") || key == "product_btn_2" {
			return "Paket 2: Paket Bisnis Premium"
		}
		if strings.HasSuffix(key, "_3") || strings.HasPrefix(key, "feature_custom_") || key == "product_btn_3" || key == "product_price_negotiable" {
			return "Paket 3: Cetak Biru & Enterprise Kustom"
		}
		return "Pengaturan Umum / Tombol"
	case 4: // Navigasi & Umum
		if strings.HasPrefix(key, "nav_") || key == "hi_user" {
			return "Menu Navigasi Header"
		}
		if strings.HasPrefix(key, "footer_") {
			return "Footer Website"
		}
		if strings.HasPrefix(key, "login_") || strings.HasPrefix(key, "register_") || strings.Contains(key, "_label") || strings.Contains(key, "_here") || key == "btn_login" || key == "btn_register" || key == "no_account" || key == "already_account" || key == "signin_here" {
			return "Halaman Login & Register"
		}
		if strings.HasPrefix(key, "chat_") {
			return "Live Chat & AI Assistant"
		}
	}
	return "Umum / Lainnya"
}
