package main

import (
	"culturstone/config"
	"culturstone/controllers"
	_ "culturstone/docs"
	"culturstone/middlewares"
	"culturstone/models"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Set environment variables for local development
	if os.Getenv("POSTGRES_URL") == "" {
		os.Setenv("POSTGRES_URL", "postgresql://neondb_owner:npg_FnMJCLO7il3b@ep-calm-fire-ahbt7y88-pooler.c-3.us-east-1.aws.neon.tech/neondb?sslmode=require")
	}

	// 1. Connect to Database (once at startup)
	config.ConnectDB()

	// 2. Initialize Gin Router
	r := gin.Default()

	// 3. Create uploads folder
	os.MkdirAll("./uploads", os.ModePerm)

	// 4. CORS Configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// 5. Visitor Counter Middleware
	r.Use(controllers.TrackVisitor)

	// 6. Static Files
	r.Static("/uploads", "./uploads")

	// 7. Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 8. API Routes
	api := r.Group("/api")
	{
		// Public Routes
		api.GET("/products", controllers.GetProducts)
		api.GET("/portfolios", controllers.GetPortfolios)
		api.GET("/testimoni", controllers.GetTestimoni)
		api.GET("/categories", controllers.GetCategories)
		api.GET("/product-categories", controllers.GetProductCategories)
		api.POST("/contact", controllers.PostContact)

		// Testimoni CRUD (Public)
		api.POST("/testimoni", controllers.SaveTestimoni)
		api.PUT("/testimoni/:id", controllers.SaveTestimoni)
		api.DELETE("/testimoni/:id", controllers.DeleteEntity(&models.Testimoni{}))

		// Auth Routes
		api.POST("/login", controllers.AdminLogin)
		api.POST("/register", controllers.AdminRegister)
	}

	// Admin Routes (Protected by JWT)
	adminGroup := api.Group("/admin")
	adminGroup.Use(middlewares.JwtAuthMiddleware())
	{
		adminGroup.GET("/stats", controllers.GetAdminStats)
		adminGroup.GET("/messages", controllers.GetMessages)
		adminGroup.DELETE("/messages/:id", controllers.DeleteEntity(&models.Message{}))
		adminGroup.GET("/visitors", controllers.GetVisitors)

		adminGroup.POST("/products", controllers.SaveProduct)
		adminGroup.PUT("/products/:id", controllers.SaveProduct)
		adminGroup.DELETE("/products/:id", controllers.DeleteEntity(&models.Product{}))

		adminGroup.POST("/categories", controllers.SaveCategory)
		adminGroup.PUT("/categories/:id", controllers.SaveCategory)
		adminGroup.DELETE("/categories/:id", controllers.DeleteEntity(&models.Category{}))

		adminGroup.POST("/portfolios", controllers.SavePortfolio)
		adminGroup.PUT("/portfolios/:id", controllers.SavePortfolio)
		adminGroup.DELETE("/portfolios/:id", controllers.DeleteEntity(&models.Portfolio{}))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Backend server running at http://localhost:%s\n", port)
	log.Fatal(r.Run(":" + port))
}
