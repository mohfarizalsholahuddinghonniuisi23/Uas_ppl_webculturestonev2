package main

import (
	"culturstone/config"
	"culturstone/controllers"
	_ "culturstone/docs"      // Dokumentasi Swagger
	"culturstone/middlewares" // Import middleware autentikasi
	"culturstone/models"
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"
)

// @title Culturestone API
// @version 1.0
// @description This is the API documentation for the Culturstone web application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	// 0. Load konfigurasi dari file .env (jika ada)
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// 1. Koneksi Database
	config.ConnectDB()

	// 1a. Auto-seed admin dari .env jika belum ada
	seedAdmin()

	// 2. Inisialisasi Router
	r := gin.Default()

	// 3. Buat folder uploads jika belum ada
	os.MkdirAll("./uploads", os.ModePerm)

	// 4. Konfigurasi CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = false // Tidak bisa AllowAllOrigins + AllowCredentials=true bersamaan
	r.Use(cors.New(corsConfig))

	// 5. Middleware Global (Visitor Counter — hanya untuk API routes)
	// Dipasang di level group, bukan global, untuk menghindari tracking static file

	// 6. Serve Static Files (Gambar/Video yang diupload lokal)
	r.Static("/uploads", "./uploads")

	// 7. Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 8. Definisi Routes
	api := r.Group("/api")
	api.Use(controllers.TrackVisitor) // Track hanya request API
	{
		// --- PUBLIC ROUTES (Tanpa Login) ---
		api.GET("/products", controllers.GetProducts)
		api.GET("/portfolios", controllers.GetPortfolios)
		api.GET("/testimoni", controllers.GetTestimoni)
		api.GET("/categories", controllers.GetCategories)
		api.GET("/product-categories", controllers.GetProductCategories)
		api.POST("/contact", controllers.PostContact)

		// Testimoni - hanya POST yang public (tambah testimoni baru boleh tanpa login)
		api.POST("/testimoni", controllers.SaveTestimoni)

		// Auth Routes (Login tetap Public, Register dihapus dari publik)
		api.POST("/login", controllers.AdminLogin)
		// CATATAN: /register sengaja dihapus dari route publik.
		// Admin dibuat otomatis via seedAdmin() saat server start.
	}

	// --- ADMIN ROUTES (DILINDUNGI JWT) ---
	adminGroup := api.Group("/admin")
	adminGroup.Use(middlewares.JwtAuthMiddleware())
	{
		// Dashboard Stats
		adminGroup.GET("/stats", controllers.GetAdminStats)

		// Messages
		adminGroup.GET("/messages", controllers.GetMessages)
		adminGroup.DELETE("/messages/:id", controllers.DeleteEntity(&models.Message{}))

		// Visitors
		adminGroup.GET("/visitors", controllers.GetVisitors)

		// CRUD Products
		adminGroup.POST("/products", controllers.SaveProduct)
		adminGroup.PUT("/products/:id", controllers.SaveProduct)
		adminGroup.DELETE("/products/:id", controllers.DeleteEntity(&models.Product{}))

		// CRUD Categories
		adminGroup.POST("/categories", controllers.SaveCategory)
		adminGroup.PUT("/categories/:id", controllers.SaveCategory)
		adminGroup.DELETE("/categories/:id", controllers.DeleteEntity(&models.Category{}))

		// CRUD Portfolios
		adminGroup.POST("/portfolios", controllers.SavePortfolio)
		adminGroup.PUT("/portfolios/:id", controllers.SavePortfolio)
		adminGroup.DELETE("/portfolios/:id", controllers.DeleteEntity(&models.Portfolio{}))

		// CRUD Testimoni — PUT & DELETE dipindah ke admin (butuh JWT)
		adminGroup.PUT("/testimoni/:id", controllers.SaveTestimoni)
		adminGroup.DELETE("/testimoni/:id", controllers.DeleteEntity(&models.Testimoni{}))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Printf("Swagger tersedia di http://localhost:%s/swagger/index.html", port)
	r.Run(":" + port)
}

// seedAdmin membuat akun admin secara otomatis dari .env saat server pertama kali dijalankan.
// Jika admin sudah ada di database, fungsi ini akan dilewati tanpa melakukan apa-apa.
func seedAdmin() {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	// Jika .env tidak mengatur ADMIN_USERNAME / ADMIN_PASSWORD, lewati
	if username == "" || password == "" {
		log.Println("[SeedAdmin] ADMIN_USERNAME atau ADMIN_PASSWORD tidak diset di .env, melewati auto-seed.")
		return
	}

	// Cek apakah sudah ada admin di database
	var count int64
	config.DB.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		log.Println("[SeedAdmin] Admin sudah ada di database, melewati auto-seed.")
		return
	}

	// Validasi kekuatan password
	if err := validateAdminPassword(password); err != nil {
		log.Fatalf("[SeedAdmin] GAGAL: Password di .env tidak memenuhi syarat keamanan: %v\n"+
			"Ubah ADMIN_PASSWORD di file .env (min 8 karakter, huruf besar, huruf kecil, angka).", err)
	}

	// Hash password dengan bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("[SeedAdmin] GAGAL hash password: %v", err)
	}

	// Simpan admin ke database
	admin := models.Admin{
		Username: username,
		Password: string(hashed),
	}
	if err := config.DB.Create(&admin).Error; err != nil {
		log.Fatalf("[SeedAdmin] GAGAL menyimpan admin ke database: %v", err)
	}

	log.Printf("[SeedAdmin] ✅ Admin '%s' berhasil dibuat otomatis dari .env.", username)
}

// validateAdminPassword memvalidasi kekuatan password.
func validateAdminPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password minimal 8 karakter")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}
	if !hasUpper {
		return fmt.Errorf("harus ada huruf kapital (A-Z)")
	}
	if !hasLower {
		return fmt.Errorf("harus ada huruf kecil (a-z)")
	}
	if !hasDigit {
		return fmt.Errorf("harus ada angka (0-9)")
	}
	return nil
}
