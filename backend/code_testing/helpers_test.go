package tests

import (
	"bytes"
	"culturstone/config"
	"culturstone/controllers"
	"culturstone/middlewares"
	"culturstone/models"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupTestDB initializes an in-memory SQLite database for testing.
// Each test gets a fresh, isolated database.
func setupTestDB(t *testing.T) {
	t.Helper()
	var err error
	config.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = config.DB.AutoMigrate(
		&models.Admin{},
		&models.Category{},
		&models.ProductCategory{},
		&models.Product{},
		&models.Portfolio{},
		&models.Testimoni{},
		&models.Message{},
		&models.Visitor{},
		&models.User{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
}

// setupRouter creates a full Gin router with all routes configured,
// mirroring the production route setup.
func setupRouter() *gin.Engine {
	r := gin.Default()

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

	return r
}

// registerAndLogin registers an admin account and returns the JWT token.
func registerAndLogin(t *testing.T, router *gin.Engine, username, password string) string {
	t.Helper()

	// Register
	regBody := map[string]string{"username": username, "password": password}
	regJSON, _ := json.Marshal(regBody)
	regReq := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(regJSON))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	router.ServeHTTP(regW, regReq)
	if regW.Code != http.StatusCreated {
		t.Fatalf("Helper: Registration failed: status %d, body: %s", regW.Code, regW.Body.String())
	}

	// Login
	loginBody := map[string]string{"username": username, "password": password}
	loginJSON, _ := json.Marshal(loginBody)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(loginJSON))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)
	if loginW.Code != http.StatusOK {
		t.Fatalf("Helper: Login failed: status %d, body: %s", loginW.Code, loginW.Body.String())
	}

	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatal("Helper: No token received from login")
	}
	return token
}

// createMultipartRequest creates a multipart form request with the given fields.
func createMultipartRequest(t *testing.T, method, url string, fields map[string]string, token string) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range fields {
		writer.WriteField(key, val)
	}
	writer.Close()

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req, httptest.NewRecorder()
}
