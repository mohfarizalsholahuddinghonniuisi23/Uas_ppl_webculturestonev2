package controllers

import (
	"culturstone/config"
	"culturstone/models"
	"net/http"
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

func PostContact(c *gin.Context) {
	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err == nil {
		config.DB.Create(&msg)
		c.JSON(http.StatusOK, gin.H{"status": "Terkirim"})
	} else {
		// Fallback form data
		msg.Name = c.PostForm("name")
		msg.Phone = c.PostForm("phone")
		msg.Email = c.PostForm("email")
		msg.Message = c.PostForm("message")
		config.DB.Create(&msg)
		c.JSON(http.StatusOK, gin.H{"status": "Terkirim"})
	}
}
