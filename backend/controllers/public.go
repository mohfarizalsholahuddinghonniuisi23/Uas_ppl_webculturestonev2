package controllers

import (
	"culturstone/config"
	"culturstone/models"
	"html"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware Visitor Counter
func TrackVisitor(c *gin.Context) {
	// Track all GET requests
	if c.Request.Method == "GET" {
		today := time.Now().Format("2006-01-02")
		var v models.Visitor
		if result := config.DB.Where("date = ?", today).First(&v); result.Error != nil {
			config.DB.Create(&models.Visitor{Date: today, Count: 1})
		} else {
			v.Count++
			config.DB.Save(&v)
		}
	}
	c.Next()
}

// GetProducts godoc
// @Summary Get all products
// @Description Get a list of all products, with an optional search query.
// @Tags Public
// @Produce json
// @Param search query string false "Search query to filter products by name"
// @Success 200 {array} models.Product
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var items []models.Product
	query := config.DB.Order("id desc")
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	query.Find(&items)
	c.JSON(http.StatusOK, items)
}

func GetPortfolios(c *gin.Context) {
	var items []models.Portfolio
	config.DB.Order("id desc").Find(&items)
	c.JSON(http.StatusOK, items)
}

// GetTestimoni for public display
func GetTestimoni(c *gin.Context) {
	var items []models.Testimoni
	config.DB.Order("id desc").Preload("Portfolio").Find(&items)

	// If no testimoni found, return empty array (not fallback to portfolios)
	if items == nil {
		items = []models.Testimoni{}
	}

	c.JSON(http.StatusOK, items)
}

func GetCategories(c *gin.Context) {
	var items []models.Category
	config.DB.Order("name asc").Find(&items)
	c.JSON(http.StatusOK, items)
}

func GetProductCategories(c *gin.Context) {
	var items []models.ProductCategory
	config.DB.Order("name asc").Find(&items)
	c.JSON(http.StatusOK, items)
}

// regexPhoneOnly hanya mengizinkan angka, +, -, dan spasi untuk nomor HP.
// BUG-06 FIX: Validasi format nomor telepon.
var regexPhoneOnly = regexp.MustCompile(`^[0-9+\-\s]{8,20}$`)

// regexNameOnly hanya mengizinkan huruf (termasuk huruf beraksara) dan spasi untuk nama.
// BUG-07 FIX: Validasi format nama (tidak boleh mengandung angka/karakter aneh).
var regexNameOnly = regexp.MustCompile(`^[\p{L}\s]{2,100}$`)

func PostContact(c *gin.Context) {
	var msg models.Message

	// Coba bind sebagai JSON terlebih dahulu
	if err := c.ShouldBindJSON(&msg); err != nil {
		// Fallback: coba form data
		msg.Name = c.PostForm("name")
		msg.Phone = c.PostForm("phone")
		msg.Email = c.PostForm("email")
		msg.Message = c.PostForm("message")
	}

	// Validasi field wajib
	if msg.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama tidak boleh kosong"})
		return
	}
	if msg.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak boleh kosong"})
		return
	}
	if msg.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesan tidak boleh kosong"})
		return
	}

	// BUG-07 FIX: Validasi nama hanya boleh huruf dan spasi
	if !regexNameOnly.MatchString(strings.TrimSpace(msg.Name)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama hanya boleh mengandung huruf dan spasi (2-100 karakter)"})
		return
	}

	// BUG-06 FIX: Validasi nomor HP hanya boleh angka, +, -, spasi
	if msg.Phone != "" && !regexPhoneOnly.MatchString(strings.TrimSpace(msg.Phone)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor HP tidak valid. Hanya angka, +, - yang diperbolehkan (8-20 karakter)"})
		return
	}

	// BUG-04 FIX: Sanitasi input menggunakan html.EscapeString untuk mencegah XSS.
	// Meskipun React sudah escape output secara default, sanitasi di backend
	// melindungi jika data diakses via API langsung atau client lain.
	msg.Name = html.EscapeString(strings.TrimSpace(msg.Name))
	msg.Phone = html.EscapeString(strings.TrimSpace(msg.Phone))
	msg.Email = html.EscapeString(strings.TrimSpace(msg.Email))
	msg.Message = html.EscapeString(strings.TrimSpace(msg.Message))

	if err := config.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pesan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Pesan berhasil terkirim"})
}
