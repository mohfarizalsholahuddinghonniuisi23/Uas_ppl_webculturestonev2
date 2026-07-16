package controllers

import (
	"culturstone/config"
	"culturstone/models"
	"culturstone/services"
	"culturstone/utils"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// BUG-05 FIX: validatePassword memvalidasi kekuatan password admin.
// Password harus ≥8 karakter, mengandung huruf besar, huruf kecil, dan angka.
func validatePassword(password string) error {
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
		return fmt.Errorf("password harus mengandung minimal 1 huruf kapital")
	}
	if !hasLower {
		return fmt.Errorf("password harus mengandung minimal 1 huruf kecil")
	}
	if !hasDigit {
		return fmt.Errorf("password harus mengandung minimal 1 angka")
	}
	return nil
}

type AdminInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AdminRegister godoc
// @Summary Register a new admin
// @Description Creates a new administrator account. Hanya bisa digunakan jika belum ada admin.
// @Tags Public
// @Accept  json
// @Produce  json
// @Param   admin  body   AdminInput  true  "Admin Credentials"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Router /register [post]
func AdminRegister(c *gin.Context) {
	// BUG-01 FIX: Cek apakah sudah ada admin di database.
	// Jika sudah ada, tolak pendaftaran (hanya 1 admin diizinkan via endpoint publik ini).
	var count int64
	config.DB.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Registrasi admin tidak diizinkan. Admin sudah terdaftar."})
		return
	}

	var input AdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// BUG-05 FIX: Validasi kekuatan password
	if err := validatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	admin := models.Admin{Username: input.Username, Password: string(hashed)}
	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// BUG-02 FIX: Response tidak mengirim objek admin (yang mengandung hash password).
	// Hanya kirim message dan username.
	c.JSON(http.StatusCreated, gin.H{"message": "Admin berhasil didaftarkan", "username": admin.Username})
}

// AdminLogin godoc
// @Summary Login as an admin
// @Description Authenticates an administrator and returns a token.
// @Tags Public
// @Accept  json
// @Produce  json
// @Param   admin  body   AdminInput  true  "Admin Credentials"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /login [post]
func AdminLogin(c *gin.Context) {
	var input AdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	// Cari admin berdasarkan username
	if err := config.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// --- GENERATE JWT TOKEN ---
	token, err := utils.GenerateToken(admin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Kembalikan token ke frontend
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"token":    token,
		"username": admin.Username,
	})
}

// GetAdminStats godoc
// @Summary Get dashboard statistics
// @Description Retrieves a summary of website data for the admin dashboard.
// @Tags Admin
// @Produce  json
// @Success 200 {object} gin.H
// @Router /admin/stats [get]
func GetAdminStats(c *gin.Context) {
	var tProd, tMsg, tPort, tTestimoni, tCat, tPCat int64
	var visitors []models.Visitor
	config.DB.Model(&models.Product{}).Count(&tProd)
	config.DB.Model(&models.Message{}).Count(&tMsg)
	config.DB.Model(&models.Portfolio{}).Count(&tPort)
	config.DB.Model(&models.Testimoni{}).Count(&tTestimoni)
	config.DB.Model(&models.Category{}).Count(&tCat)
	config.DB.Model(&models.ProductCategory{}).Count(&tPCat)
	config.DB.Order("date desc").Limit(7).Find(&visitors)
	c.JSON(http.StatusOK, gin.H{
		"products":           tProd,
		"messages":           tMsg,
		"portfolios":         tPort,
		"testimoni":          tTestimoni,
		"categories":         tCat,
		"product_categories": tPCat,
		"visitors":           visitors,
	})
}

// GetMessages godoc
// @Summary Get all contact messages
// @Description Retrieves a list of all messages submitted through the contact form.
// @Tags Admin
// @Produces  json
// @Success 200 {array} models.Message
// @Router /admin/messages [get]
func GetMessages(c *gin.Context) {
	var msgs []models.Message
	config.DB.Order("created_at desc").Find(&msgs)
	c.JSON(http.StatusOK, msgs)
}

// SaveProduct godoc
// @Summary Create or update a product
// @Description Creates a new product or updates an existing one if an ID is provided.
// @Tags Admin
// @Consumes  multipart/form-data
// @Produces  json
// @Param id path int false "Product ID"
// @Param name formData string true "Product Name"
// @Param description formData string true "Product Description"
// @Param category formData string true "Product Category"
// @Param quality formData string false "Product Quality"
// @Param size formData string false "Product Size"
// @Param finishing formData string false "Product Finishing"
// @Param image formData file false "Product Image"
// @Success 200 {object} models.Product
// @Router /admin/products [post]
// @Router /admin/products/{id} [put]
func SaveProduct(c *gin.Context) {
	var p models.Product
	id := c.Param("id")
	if id != "" {
		config.DB.First(&p, id)
	}

	p.Name = c.PostForm("name")
	p.Description = c.PostForm("description")
	p.Category = c.PostForm("category")
	p.Quality = c.PostForm("quality")
	p.Size = c.PostForm("size")
	p.Finishing = c.PostForm("finishing")

	file, err := c.FormFile("image")
	if err == nil { // If a file was provided
		// Upload file to Vercel Blob
		imageURL, uploadErr := services.UploadFileToVercelBlob(file, "prod")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Blob: " + uploadErr.Error()})
			return
		}
		p.ImageURL = imageURL
	} else if err != http.ErrMissingFile {
		// Handle other errors from c.FormFile besides just a missing file (which is fine if not updating image)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file: " + err.Error()})
		return
	}

	// Log start of database operation
	startDBOp := time.Now()
	if id != "" {
		config.DB.Save(&p)
		log.Printf("Product ID %s updated in %v", id, time.Since(startDBOp))
	} else {
		config.DB.Create(&p)
		log.Printf("New product created in %v", time.Since(startDBOp))
	}
	c.JSON(http.StatusOK, p)
}

// SaveCategory godoc
// @Summary Create or update a category
// @Description Creates a new category or updates an existing one if an ID is provided.
// @Tags Admin
// @Consumes  multipart/form-data
// @Produces  json
// @Param id path int false "Category ID"
// @Param name formData string true "Category Name"
// @Param description formData string false "Category Description"
// @Param image formData file false "Category Image"
// @Success 200 {object} models.Category
// @Router /admin/categories [post]
// @Router /admin/categories/{id} [put]
func SaveCategory(c *gin.Context) {
	var cat models.Category
	id := c.Param("id")
	if id != "" {
		config.DB.First(&cat, id)
	}

	cat.Name = c.PostForm("name")
	cat.Description = c.PostForm("description")

	file, err := c.FormFile("image")
	if err == nil { // If a file was provided
		// Upload file to Vercel Blob
		imageURL, uploadErr := services.UploadFileToVercelBlob(file, "cat")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Blob: " + uploadErr.Error()})
			return
		}
		cat.ImageURL = imageURL
	} else if err != http.ErrMissingFile {
		// Handle other errors from c.FormFile besides just a missing file (which is fine if not updating image)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file: " + err.Error()})
		return
	}

	// Log start of database operation
	startDBOp := time.Now()
	if id != "" {
		config.DB.Save(&cat)
		log.Printf("Category ID %s updated in %v", id, time.Since(startDBOp))
	} else {
		config.DB.Create(&cat)
		log.Printf("New category created in %v", time.Since(startDBOp))
	}
	c.JSON(http.StatusOK, cat)
}

// SavePortfolio godoc
// @Summary Create or update a portfolio item
// @Description Creates a new portfolio item or updates an existing one if an ID is provided.
// @Tags Admin
// @Consumes  multipart/form-data
// @Produces  json
// @Param id path int false "Portfolio ID"
// @Param title formData string true "Portfolio Title"
// @Param description formData string false "Portfolio Description"
// @Param image formData file false "Portfolio Image"
// @Success 200 {object} models.Portfolio
// @Router /admin/portfolios [post]
// @Router /admin/portfolios/{id} [put]
func SavePortfolio(c *gin.Context) {
	var pf models.Portfolio
	id := c.Param("id")
	if id != "" {
		config.DB.First(&pf, id)
	}

	pf.Title = c.PostForm("title")
	pf.Description = c.PostForm("description")

	file, err := c.FormFile("image")
	if err == nil { // If a file was provided
		// Upload file to Vercel Blob
		imageURL, uploadErr := services.UploadFileToVercelBlob(file, "port")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Blob: " + uploadErr.Error()})
			return
		}
		pf.ImageURL = imageURL
	} else if err != http.ErrMissingFile {
		// Handle other errors from c.FormFile besides just a missing file (which is fine if not updating image)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file: " + err.Error()})
		return
	}

	// Log start of database operation
	startDBOp := time.Now()
	if id != "" {
		config.DB.Save(&pf)
		log.Printf("Portfolio ID %s updated in %v", id, time.Since(startDBOp))
	} else {
		config.DB.Create(&pf)
		log.Printf("New portfolio created in %v", time.Since(startDBOp))
	}
	c.JSON(http.StatusOK, pf)
}

// SaveTestimoni godoc
// @Summary Create or update a testimonial
// @Description Creates a new testimonial or updates an existing one if an ID is provided.
// @Tags Admin
// @Consumes  multipart/form-data
// @Produces  json
// @Param id path int false "Testimoni ID"
// @Param client_name formData string true "Client Name"
// @Param testimonial_text formData string true "Testimonial Text"
// @Param portfolio_id formData int false "Associated Portfolio ID"
// @Param image formData file false "Client Image"
// @Param video formData file false "Testimonial Video"
// @Success 200 {object} models.Testimoni
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/testimoni [post]
// @Router /admin/testimoni/{id} [put]
// SaveTestimoni godoc
func SaveTestimoni(c *gin.Context) {
	var t models.Testimoni
	id := c.Param("id")

	// 1. Cek Data Lama (Edit Mode)
	if id != "" {
		if err := config.DB.First(&t, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Testimoni not found"})
			return
		}
	}

	// 2. Ambil Input Text
	t.ClientName = c.PostForm("client_name")
	t.TestimonialText = c.PostForm("testimonial_text")

	// 3. LOGIC BARU: Handle Portfolio ID agar bisa NULL
	portfolioIDStr := c.PostForm("portfolio_id")

	// SELALU set PortfolioID ke nil terlebih dahulu
	t.PortfolioID = nil

	if portfolioIDStr != "" && portfolioIDStr != "0" {
		// Jika user memilih portfolio
		pID, err := strconv.ParseUint(portfolioIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid portfolio_id format"})
			return
		}

		// Validasi portfolio ada di database
		uID := uint(pID)
		var portfolio models.Portfolio
		if err := config.DB.First(&portfolio, uID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Portfolio with ID " + portfolioIDStr + " not found"})
			return
		}

		t.PortfolioID = &uID
	}

	// 4. Upload Image (Tetap sama)
	file, err := c.FormFile("image")
	if err == nil {
		imageURL, uploadErr := services.UploadFileToVercelBlob(file, "testimoni_img")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + uploadErr.Error()})
			return
		}
		t.ImagePath = imageURL
	}

	// 5. Upload Video (Tetap sama)
	videoFile, err := c.FormFile("video")
	if err == nil {
		videoURL, uploadErr := services.UploadFileToVercelBlob(videoFile, "testimoni_video")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video: " + uploadErr.Error()})
			return
		}
		t.VideoPath = videoURL
	}

	// 6. Simpan ke Database
	startDBOp := time.Now()
	var errSave error

	if id != "" {
		// Update dengan Map agar NULL bisa tersimpan dengan benar
		updateData := map[string]interface{}{
			"client_name":      t.ClientName,
			"testimonial_text": t.TestimonialText,
			"image_path":       t.ImagePath,
			"video_path":       t.VideoPath,
		}
		// Handle portfolio_id: jika nil, set ke NULL
		if t.PortfolioID == nil {
			updateData["portfolio_id"] = nil
		} else {
			updateData["portfolio_id"] = *t.PortfolioID
		}
		errSave = config.DB.Model(&t).Updates(updateData).Error
	} else {
		// Create Baru - set portfolio_id ke NULL jika nil
		if t.PortfolioID == nil {
			t.PortfolioID = nil // Pastikan nil
		}
		errSave = config.DB.Create(&t).Error
	}

	if errSave != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create testimoni: " + errSave.Error()})
		return
	}

	log.Printf("Testimoni saved successfully in %v", time.Since(startDBOp))
	c.JSON(http.StatusOK, t)
}

// AdminGetTestimoni godoc
// @Summary Get all testimonials for admin
// @Description Retrieves a list of all testimonials.
// @Tags Admin
// @Produces json
// @Success 200 {array} models.Testimoni
// @Router /admin/testimoni [get]
func AdminGetTestimoni(c *gin.Context) {
	var testimoni []models.Testimoni
	config.DB.Order("id desc").Preload("Portfolio").Find(&testimoni)
	c.JSON(http.StatusOK, testimoni)
}

// Generic Delete
// Menggunakan reflect untuk membuat instance baru setiap request,
// sehingga tidak ada shared state antar request (menghindari race condition).
func DeleteEntity(model interface{}) gin.HandlerFunc {
	// Ambil tipe dari model yang diberikan (misal: *models.Product → models.Product)
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak boleh kosong"})
			return
		}

		// Buat instance baru dari tipe model setiap request
		newInstance := reflect.New(modelType).Interface()

		result := config.DB.Delete(newInstance, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data: " + result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data dengan ID " + id + " tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "deleted", "id": id})
	}
}

// GetVisitors godoc
// @Summary Get all visitor records
// @Description Retrieves a list of daily visitor counts.
// @Tags Admin
// @Produces json
// @Success 200 {array} models.Visitor
// @Router /admin/visitors [get]
func GetVisitors(c *gin.Context) {
	var visitors []models.Visitor
	config.DB.Order("date desc").Find(&visitors)
	c.JSON(http.StatusOK, visitors)
}
