package tests

// ============================================================================
// UNIT TESTING — Aplikasi Web Culturstone
// ============================================================================
// Unit test menguji setiap fungsi/method secara individual dan terisolasi.
// Database menggunakan SQLite in-memory agar setiap test independen.
// Setiap fitur diuji langsung pada level controller function tanpa melalui router.
// ============================================================================

import (
	"bytes"
	"culturstone/config"
	"culturstone/controllers"
	"culturstone/models"
	"culturstone/utils"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// ============================================================
// FITUR 1: Register Admin
// ============================================================

func TestUnit_F01_RegisterAdmin_ValidData(t *testing.T) {
	setupTestDB(t)
	body, _ := json.Marshal(controllers.AdminInput{Username: "admin1", Password: "pass123"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controllers.AdminRegister(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201, got %d. Body: %s", w.Code, w.Body.String())
	}
	t.Log("✅ F01-UT01: Register admin dengan data valid berhasil")
}

func TestUnit_F01_RegisterAdmin_EmptyBody(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/register", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	controllers.AdminRegister(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", w.Code)
	}
	t.Log("✅ F01-UT02: Register tanpa body ditolak (400)")
}

func TestUnit_F01_RegisterAdmin_DuplicateUsername(t *testing.T) {
	setupTestDB(t)

	// Register pertama
	body1, _ := json.Marshal(controllers.AdminInput{Username: "dupuser", Password: "pass1"})
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body1))
	c1.Request.Header.Set("Content-Type", "application/json")
	controllers.AdminRegister(c1)

	// Register kedua (duplikat)
	body2, _ := json.Marshal(controllers.AdminInput{Username: "dupuser", Password: "pass2"})
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body2))
	c2.Request.Header.Set("Content-Type", "application/json")
	controllers.AdminRegister(c2)

	if w2.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 for duplicate, got %d", w2.Code)
	}
	t.Log("✅ F01-UT03: Register username duplikat ditolak (400)")
}

// ============================================================
// FITUR 2: Login Admin
// ============================================================

func TestUnit_F02_LoginAdmin_ValidCredentials(t *testing.T) {
	setupTestDB(t)

	// Register dulu
	regBody, _ := json.Marshal(controllers.AdminInput{Username: "loginuser", Password: "loginpass"})
	rw := httptest.NewRecorder()
	rc, _ := gin.CreateTestContext(rw)
	rc.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(regBody))
	rc.Request.Header.Set("Content-Type", "application/json")
	controllers.AdminRegister(rc)

	// Login
	loginBody, _ := json.Marshal(controllers.AdminInput{Username: "loginuser", Password: "loginpass"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(loginBody))
	c.Request.Header.Set("Content-Type", "application/json")
	controllers.AdminLogin(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == nil || resp["token"] == "" {
		t.Fatal("Expected token in response")
	}
	t.Log("✅ F02-UT01: Login dengan kredensial valid berhasil, token diterima")
}

func TestUnit_F02_LoginAdmin_WrongUsername(t *testing.T) {
	setupTestDB(t)
	body, _ := json.Marshal(controllers.AdminInput{Username: "ghost", Password: "pass"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	controllers.AdminLogin(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", w.Code)
	}
	t.Log("✅ F02-UT02: Login username salah ditolak (401)")
}

func TestUnit_F02_LoginAdmin_WrongPassword(t *testing.T) {
	setupTestDB(t)

	// Register
	regBody, _ := json.Marshal(controllers.AdminInput{Username: "pwuser", Password: "correct"})
	rw := httptest.NewRecorder()
	rc, _ := gin.CreateTestContext(rw)
	rc.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(regBody))
	rc.Request.Header.Set("Content-Type", "application/json")
	controllers.AdminRegister(rc)

	// Login dengan password salah
	loginBody, _ := json.Marshal(controllers.AdminInput{Username: "pwuser", Password: "wrong"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(loginBody))
	c.Request.Header.Set("Content-Type", "application/json")
	controllers.AdminLogin(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", w.Code)
	}
	t.Log("✅ F02-UT03: Login password salah ditolak (401)")
}

func TestUnit_F02_LoginAdmin_EmptyBody(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	controllers.AdminLogin(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", w.Code)
	}
	t.Log("✅ F02-UT04: Login tanpa body ditolak (400)")
}

// ============================================================
// FITUR 3: Get Products (+ Search)
// ============================================================

func TestUnit_F03_GetProducts_All(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Product{Name: "Granit Hitam", Category: "Eksterior"})
	config.DB.Create(&models.Product{Name: "Marmer Putih", Category: "Dekorasi"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/products", nil)
	controllers.GetProducts(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var products []models.Product
	json.Unmarshal(w.Body.Bytes(), &products)
	if len(products) != 2 {
		t.Fatalf("Expected 2 products, got %d", len(products))
	}
	t.Log("✅ F03-UT01: GetProducts mengembalikan semua produk")
}

func TestUnit_F03_GetProducts_WithSearch(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Product{Name: "Granit Hitam"})
	config.DB.Create(&models.Product{Name: "Marmer Putih"})
	config.DB.Create(&models.Product{Name: "Granit Abu"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/products?search=Granit", nil)
	controllers.GetProducts(c)

	var products []models.Product
	json.Unmarshal(w.Body.Bytes(), &products)
	if len(products) != 2 {
		t.Fatalf("Expected 2 products matching 'Granit', got %d", len(products))
	}
	t.Log("✅ F03-UT02: GetProducts dengan search filter berfungsi")
}

func TestUnit_F03_GetProducts_SearchNoResults(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Product{Name: "Granit Hitam"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/products?search=XXXXXXXX", nil)
	controllers.GetProducts(c)

	var products []models.Product
	json.Unmarshal(w.Body.Bytes(), &products)
	if len(products) != 0 {
		t.Fatalf("Expected 0 products, got %d", len(products))
	}
	t.Log("✅ F03-UT03: GetProducts search tanpa hasil mengembalikan array kosong")
}

// ============================================================
// FITUR 4: Get Portfolios
// ============================================================

func TestUnit_F04_GetPortfolios(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Portfolio{Title: "Proyek Hotel", Description: "Instalasi marmer"})
	config.DB.Create(&models.Portfolio{Title: "Proyek Masjid", Description: "Instalasi granit"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/portfolios", nil)
	controllers.GetPortfolios(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var items []models.Portfolio
	json.Unmarshal(w.Body.Bytes(), &items)
	if len(items) != 2 {
		t.Fatalf("Expected 2 portfolios, got %d", len(items))
	}
	t.Log("✅ F04-UT01: GetPortfolios mengembalikan semua portofolio")
}

// ============================================================
// FITUR 5: Get Testimoni
// ============================================================

func TestUnit_F05_GetTestimoni_Empty(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/testimoni", nil)
	controllers.GetTestimoni(c)

	if w.Body.String() != "[]" {
		t.Fatalf("Expected '[]', got: %s", w.Body.String())
	}
	t.Log("✅ F05-UT01: GetTestimoni data kosong mengembalikan array kosong []")
}

func TestUnit_F05_GetTestimoni_WithData(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Testimoni{ClientName: "Budi", TestimonialText: "Bagus!"})
	config.DB.Create(&models.Testimoni{ClientName: "Andi", TestimonialText: "Mantap!"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/testimoni", nil)
	controllers.GetTestimoni(c)

	var items []models.Testimoni
	json.Unmarshal(w.Body.Bytes(), &items)
	if len(items) != 2 {
		t.Fatalf("Expected 2 testimoni, got %d", len(items))
	}
	t.Log("✅ F05-UT02: GetTestimoni mengembalikan data testimoni")
}

// ============================================================
// FITUR 6: Get Categories
// ============================================================

func TestUnit_F06_GetCategories(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Category{Name: "Granit", Description: "Batu keras"})
	config.DB.Create(&models.Category{Name: "Marmer", Description: "Mewah"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	controllers.GetCategories(c)

	var items []models.Category
	json.Unmarshal(w.Body.Bytes(), &items)
	if len(items) != 2 {
		t.Fatalf("Expected 2 categories, got %d", len(items))
	}
	t.Log("✅ F06-UT01: GetCategories mengembalikan semua kategori")
}

// ============================================================
// FITUR 7: Get Product Categories
// ============================================================

func TestUnit_F07_GetProductCategories(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.ProductCategory{Name: "Sanitary"})
	config.DB.Create(&models.ProductCategory{Name: "Dekorasi"})
	config.DB.Create(&models.ProductCategory{Name: "Eksterior"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/product-categories", nil)
	controllers.GetProductCategories(c)

	var items []models.ProductCategory
	json.Unmarshal(w.Body.Bytes(), &items)
	if len(items) != 3 {
		t.Fatalf("Expected 3 product-categories, got %d", len(items))
	}
	t.Log("✅ F07-UT01: GetProductCategories mengembalikan semua kategori produk")
}

// ============================================================
// FITUR 8: Post Contact Form
// ============================================================

func TestUnit_F08_PostContact_ValidJSON(t *testing.T) {
	setupTestDB(t)
	msg := map[string]string{"Name": "Budi", "Phone": "081234567890", "Email": "budi@mail.com", "Message": "Halo"}
	body, _ := json.Marshal(msg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/contact", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	controllers.PostContact(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "Terkirim" {
		t.Fatalf("Expected 'Terkirim', got: %v", resp["status"])
	}
	// Verifikasi tersimpan di DB
	var messages []models.Message
	config.DB.Find(&messages)
	if len(messages) != 1 {
		t.Fatalf("Expected 1 message in DB, got %d", len(messages))
	}
	t.Log("✅ F08-UT01: PostContact JSON valid berhasil tersimpan")
}

func TestUnit_F08_PostContact_MalformedJSON(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/contact", bytes.NewBufferString("bukan json"))
	c.Request.Header.Set("Content-Type", "application/json")
	controllers.PostContact(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "Terkirim" {
		t.Fatalf("Expected fallback 'Terkirim', got: %v", resp["status"])
	}
	t.Log("✅ F08-UT02: PostContact malformed JSON ditangani via fallback")
}

// ============================================================
// FITUR 9: Save Product (Create/Update)
// ============================================================

func TestUnit_F09_SaveProduct_Create(t *testing.T) {
	setupTestDB(t)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Wastafel Marmer")
	writer.WriteField("description", "Wastafel premium")
	writer.WriteField("category", "Sanitary")
	writer.WriteField("quality", "Export")
	writer.WriteField("size", "40cm")
	writer.WriteField("finishing", "Polish")
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	controllers.SaveProduct(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	// Verifikasi di DB
	var products []models.Product
	config.DB.Find(&products)
	if len(products) != 1 || products[0].Name != "Wastafel Marmer" {
		t.Fatalf("Product not saved correctly")
	}
	t.Log("✅ F09-UT01: SaveProduct membuat produk baru berhasil")
}

// ============================================================
// FITUR 10: Save Category (Create/Update)
// ============================================================

func TestUnit_F10_SaveCategory_Create(t *testing.T) {
	setupTestDB(t)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Granit Premium")
	writer.WriteField("description", "Koleksi granit terbaik")
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	controllers.SaveCategory(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var cats []models.Category
	config.DB.Find(&cats)
	if len(cats) != 1 || cats[0].Name != "Granit Premium" {
		t.Fatal("Category not saved correctly")
	}
	t.Log("✅ F10-UT01: SaveCategory membuat kategori baru berhasil")
}

// ============================================================
// FITUR 11: Save Portfolio (Create/Update)
// ============================================================

func TestUnit_F11_SavePortfolio_Create(t *testing.T) {
	setupTestDB(t)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Proyek Hotel Bintang 5")
	writer.WriteField("description", "Pemasangan marmer lobby")
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	controllers.SavePortfolio(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var items []models.Portfolio
	config.DB.Find(&items)
	if len(items) != 1 || items[0].Title != "Proyek Hotel Bintang 5" {
		t.Fatal("Portfolio not saved correctly")
	}
	t.Log("✅ F11-UT01: SavePortfolio membuat portofolio baru berhasil")
}

// ============================================================
// FITUR 12: Save Testimoni (Create/Update)
// ============================================================

func TestUnit_F12_SaveTestimoni_Create(t *testing.T) {
	setupTestDB(t)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("client_name", "Pak Budi")
	writer.WriteField("testimonial_text", "Kualitas batu alam sangat bagus!")
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	controllers.SaveTestimoni(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	var items []models.Testimoni
	config.DB.Find(&items)
	if len(items) != 1 || items[0].ClientName != "Pak Budi" {
		t.Fatal("Testimoni not saved correctly")
	}
	t.Log("✅ F12-UT01: SaveTestimoni membuat testimoni baru berhasil")
}

func TestUnit_F12_SaveTestimoni_WithPortfolio(t *testing.T) {
	setupTestDB(t)

	// Buat portfolio terlebih dahulu
	config.DB.Create(&models.Portfolio{Title: "Test Portfolio"})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("client_name", "Pak Andi")
	writer.WriteField("testimonial_text", "Proyek hotel sangat memuaskan")
	writer.WriteField("portfolio_id", "1")
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	controllers.SaveTestimoni(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	var items []models.Testimoni
	config.DB.Find(&items)
	if items[0].PortfolioID == nil || *items[0].PortfolioID != 1 {
		t.Fatal("Testimoni portfolio_id not saved correctly")
	}
	t.Log("✅ F12-UT02: SaveTestimoni dengan portfolio_id berhasil")
}

// ============================================================
// FITUR 13: Get Admin Stats (Dashboard)
// ============================================================

func TestUnit_F13_GetAdminStats(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Product{Name: "P1"})
	config.DB.Create(&models.Product{Name: "P2"})
	config.DB.Create(&models.Message{Name: "M1"})
	config.DB.Create(&models.Portfolio{Title: "Port1"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	controllers.GetAdminStats(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var stats map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &stats)

	expectedKeys := []string{"products", "messages", "portfolios", "testimoni", "categories", "product_categories", "visitors"}
	for _, key := range expectedKeys {
		if _, exists := stats[key]; !exists {
			t.Fatalf("Missing key '%s' in stats", key)
		}
	}
	if stats["products"].(float64) != 2 {
		t.Fatalf("Expected 2 products in stats, got %v", stats["products"])
	}
	t.Log("✅ F13-UT01: GetAdminStats mengembalikan statistik lengkap dan benar")
}

// ============================================================
// FITUR 14: Get/Delete Messages
// ============================================================

func TestUnit_F14_GetMessages(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Message{Name: "User1", Email: "u1@mail.com", Message: "Halo"})
	config.DB.Create(&models.Message{Name: "User2", Email: "u2@mail.com", Message: "Hi"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	controllers.GetMessages(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var msgs []models.Message
	json.Unmarshal(w.Body.Bytes(), &msgs)
	if len(msgs) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(msgs))
	}
	t.Log("✅ F14-UT01: GetMessages mengembalikan semua pesan")
}

func TestUnit_F14_DeleteEntity_Message(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Message{Name: "ToDelete", Message: "Hapus saya"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handler := controllers.DeleteEntity(&models.Message{})
	handler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "deleted" {
		t.Fatalf("Expected 'deleted', got: %v", resp["status"])
	}
	t.Log("✅ F14-UT02: DeleteEntity(Message) berhasil menghapus pesan")
}

// ============================================================
// FITUR 15: JWT Token Utilities
// ============================================================

func TestUnit_F15_GenerateToken(t *testing.T) {
	token, err := utils.GenerateToken(1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token")
	}
	t.Log("✅ F15-UT01: GenerateToken menghasilkan token JWT valid")
}

func TestUnit_F15_ExtractToken_ValidBearer(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer mytoken123")

	result := utils.ExtractToken(c)
	if result != "mytoken123" {
		t.Fatalf("Expected 'mytoken123', got '%s'", result)
	}
	t.Log("✅ F15-UT02: ExtractToken mengekstrak token dari header Bearer")
}

func TestUnit_F15_ExtractToken_EmptyHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	result := utils.ExtractToken(c)
	if result != "" {
		t.Fatalf("Expected empty string, got '%s'", result)
	}
	t.Log("✅ F15-UT03: ExtractToken mengembalikan string kosong untuk header kosong")
}

func TestUnit_F15_TokenValid_ValidToken(t *testing.T) {
	token, _ := utils.GenerateToken(1)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	err := utils.TokenValid(c)
	if err != nil {
		t.Fatalf("Expected valid token, got error: %v", err)
	}
	t.Log("✅ F15-UT04: TokenValid memvalidasi token valid")
}

func TestUnit_F15_TokenValid_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalidtoken")

	err := utils.TokenValid(c)
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
	t.Log("✅ F15-UT05: TokenValid menolak token invalid")
}

func TestUnit_F15_TokenValid_NoHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	err := utils.TokenValid(c)
	if err == nil {
		t.Fatal("Expected error for missing header, got nil")
	}
	t.Log("✅ F15-UT06: TokenValid menolak request tanpa Authorization header")
}
