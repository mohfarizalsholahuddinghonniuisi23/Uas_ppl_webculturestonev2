package handler

import (
	"culturstone/config"
	"culturstone/controllers"
	_ "culturstone/docs"      // Dokumentasi Swagger
	"culturstone/middlewares" // Import middleware autentikasi yang baru dibuat
	"culturstone/models"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Culturstone API
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
func Handler(w http.ResponseWriter, req *http.Request) {
	// 1. Koneksi Database
	config.ConnectDB()

	// 2. Inisialisasi Router
	r := gin.Default()

	// 3. Buat folder uploads jika belum ada
	os.MkdirAll("./uploads", os.ModePerm)

	// 4. Konfigurasi CORS (Penting untuk Frontend React)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Izinkan semua origin untuk Vercel
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// 5. Middleware Global (Visitor Counter)
	r.Use(controllers.TrackVisitor)

	// 6. Serve Static Files (Gambar/Video)
	r.Static("/uploads", "./uploads")

	// 7. Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 8. Definisi Routes
	api := r.Group("/api")
	{
		// --- PUBLIC ROUTES (Tanpa Login) ---
		api.GET("/products", controllers.GetProducts)
		api.GET("/portfolios", controllers.GetPortfolios)
		api.GET("/testimoni", controllers.GetTestimoni)
		api.GET("/categories", controllers.GetCategories)
		api.GET("/product-categories", controllers.GetProductCategories)
		api.POST("/contact", controllers.PostContact)

		// Testimoni CRUD (PUBLIC - tanpa auth middleware)
		api.POST("/testimoni", controllers.SaveTestimoni)
		api.PUT("/testimoni/:id", controllers.SaveTestimoni)
		api.DELETE("/testimoni/:id", controllers.DeleteEntity(&models.Testimoni{}))

		// Auth Routes (Login & Register tetap Public)
		api.POST("/login", controllers.AdminLogin)
		api.POST("/register", controllers.AdminRegister)
	}

	// --- ADMIN ROUTES (DILINDUNGI JWT) - Defined AFTER public routes ---
	adminGroup := api.Group("/admin")

	// Terapkan Middleware Autentikasi di sini
	adminGroup.Use(middlewares.JwtAuthMiddleware())
	{
		// Dashboard Stats
		adminGroup.GET("/stats", controllers.GetAdminStats)

		// Messages
		adminGroup.GET("/messages", controllers.GetMessages)
		// @Summary Delete a contact message
		// @Router /admin/messages/{id} [delete]
		adminGroup.DELETE("/messages/:id", controllers.DeleteEntity(&models.Message{}))

		// Visitors
		adminGroup.GET("/visitors", controllers.GetVisitors)

		// CRUD Products
		adminGroup.POST("/products", controllers.SaveProduct)
		adminGroup.PUT("/products/:id", controllers.SaveProduct)
		// @Summary Delete a product
		// @Router /admin/products/{id} [delete]
		adminGroup.DELETE("/products/:id", controllers.DeleteEntity(&models.Product{}))

		// CRUD Categories
		adminGroup.POST("/categories", controllers.SaveCategory)
		adminGroup.PUT("/categories/:id", controllers.SaveCategory)
		// @Summary Delete a category
		// @Router /admin/categories/{id} [delete]
		adminGroup.DELETE("/categories/:id", controllers.DeleteEntity(&models.Category{}))

		// CRUD Portfolios
		adminGroup.POST("/portfolios", controllers.SavePortfolio)
		adminGroup.PUT("/portfolios/:id", controllers.SavePortfolio)
		// @Summary Delete a portfolio item
		// @Router /admin/portfolios/{id} [delete]
		adminGroup.DELETE("/portfolios/:id", controllers.DeleteEntity(&models.Portfolio{}))
	}

	r.ServeHTTP(w, req)
}
