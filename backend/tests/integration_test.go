package tests

// ============================================================================
// INTEGRATION TESTING — Aplikasi Web Culturstone
// ============================================================================
// Integration test memvalidasi bahwa modul-modul yang berbeda bekerja dengan
// benar ketika diintegrasikan. Pengujian menguji alur lengkap dari:
// HTTP Request → Router → Middleware → Controller → Database → Response.
// Semua 15 fitur diuji dalam skenario alur terintegrasi.
// ============================================================================

import (
	"bytes"
	"culturstone/config"
	"culturstone/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================
// FITUR 1: Register Admin — Alur Integrasi
// ============================================================

func TestIntegration_F01_RegisterFlow(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// Register admin baru
	body := map[string]string{"username": "intadmin", "password": "intpass123"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201, got %d", w.Code)
	}

	// Verifikasi di database
	var admin models.Admin
	config_err := configDBFindAdmin(&admin, "intadmin")
	if config_err {
		t.Fatal("Admin not found in database after registration")
	}
	t.Log("✅ F01-IT01: Alur register admin terintegrasi berhasil (HTTP → Controller → DB)")
}

func TestIntegration_F01_RegisterDuplicate(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body := map[string]string{"username": "dupuser", "password": "pass1"}
	jsonBody, _ := json.Marshal(body)

	// Register pertama
	req1 := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// Register kedua (duplikat)
	jsonBody2, _ := json.Marshal(body)
	req2 := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 for duplicate, got %d", w2.Code)
	}
	t.Log("✅ F01-IT02: Register duplikat ditolak melalui alur terintegrasi")
}

// ============================================================
// FITUR 2: Login Admin — Alur Register → Login
// ============================================================

func TestIntegration_F02_RegisterThenLogin(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	token := registerAndLogin(t, router, "loginadmin", "loginpass")
	if token == "" {
		t.Fatal("Expected non-empty token after register → login")
	}
	t.Log("✅ F02-IT01: Alur Register → Login berhasil, token JWT diterima")
}

func TestIntegration_F02_LoginWrongCredentials(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// Register
	registerAndLogin(t, router, "realuser", "realpass")

	// Login dengan password salah
	body := map[string]string{"username": "realuser", "password": "wrongpass"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", w.Code)
	}
	t.Log("✅ F02-IT02: Login password salah ditolak melalui router")
}

// ============================================================
// FITUR 3: Get Products — Alur Create → Get
// ============================================================

func TestIntegration_F03_CreateThenGetProducts(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	token := registerAndLogin(t, router, "prodadmin", "prodpass")

	// Create Product via admin API
	req, w := createMultipartRequest(t, http.MethodPost, "/api/admin/products",
		map[string]string{"name": "Granit Hitam", "description": "Batu granit", "category": "Eksterior"}, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Create product expected 200, got %d", w.Code)
	}

	// Get Products via public API
	getReq := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	var products []models.Product
	json.Unmarshal(getW.Body.Bytes(), &products)
	if len(products) != 1 || products[0].Name != "Granit Hitam" {
		t.Fatal("Product not found after creation")
	}
	t.Log("✅ F03-IT01: Alur Create Product → Get Products terintegrasi berhasil")
}

// ============================================================
// FITUR 4: Get Portfolios — Alur Create → Get
// ============================================================

func TestIntegration_F04_CreateThenGetPortfolios(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	token := registerAndLogin(t, router, "portadmin", "portpass")

	// Create Portfolio
	req, w := createMultipartRequest(t, http.MethodPost, "/api/admin/portfolios",
		map[string]string{"title": "Proyek Hotel", "description": "Pemasangan marmer"}, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Create portfolio expected 200, got %d", w.Code)
	}

	// Get Portfolios
	getReq := httptest.NewRequest(http.MethodGet, "/api/portfolios", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	var items []models.Portfolio
	json.Unmarshal(getW.Body.Bytes(), &items)
	if len(items) != 1 || items[0].Title != "Proyek Hotel" {
		t.Fatal("Portfolio not found after creation")
	}
	t.Log("✅ F04-IT01: Alur Create Portfolio → Get Portfolios berhasil")
}

// ============================================================
// FITUR 5: Get Testimoni — Alur Create → Get
// ============================================================

func TestIntegration_F05_CreateThenGetTestimoni(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// Create testimoni via public API (no auth needed)
	req, w := createMultipartRequest(t, http.MethodPost, "/api/testimoni",
		map[string]string{"client_name": "Budi", "testimonial_text": "Sangat puas!"}, "")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Create testimoni expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Get Testimoni
	getReq := httptest.NewRequest(http.MethodGet, "/api/testimoni", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	var items []models.Testimoni
	json.Unmarshal(getW.Body.Bytes(), &items)
	if len(items) != 1 || items[0].ClientName != "Budi" {
		t.Fatal("Testimoni not found after creation")
	}
	t.Log("✅ F05-IT01: Alur Create Testimoni → Get Testimoni berhasil")
}

// ============================================================
// FITUR 6: Get Categories — Alur Create → Get
// ============================================================

func TestIntegration_F06_CreateThenGetCategories(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	token := registerAndLogin(t, router, "catadmin", "catpass")

	// Create Category
	req, w := createMultipartRequest(t, http.MethodPost, "/api/admin/categories",
		map[string]string{"name": "Granit", "description": "Batu keras"}, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Create category expected 200, got %d", w.Code)
	}

	// Get Categories
	getReq := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	var items []models.Category
	json.Unmarshal(getW.Body.Bytes(), &items)
	if len(items) != 1 || items[0].Name != "Granit" {
		t.Fatal("Category not found after creation")
	}
	t.Log("✅ F06-IT01: Alur Create Category → Get Categories berhasil")
}

// ============================================================
// FITUR 7: Get Product Categories
// ============================================================

func TestIntegration_F07_GetProductCategories(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// Seed data langsung ke DB
	config_db_create(&models.ProductCategory{Name: "Sanitary"})
	config_db_create(&models.ProductCategory{Name: "Dekorasi"})

	req := httptest.NewRequest(http.MethodGet, "/api/product-categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var items []models.ProductCategory
	json.Unmarshal(w.Body.Bytes(), &items)
	if len(items) != 2 {
		t.Fatalf("Expected 2, got %d", len(items))
	}
	t.Log("✅ F07-IT01: Get Product Categories melalui router berhasil")
}

// ============================================================
// FITUR 8: Post Contact → Admin View Messages
// ============================================================

func TestIntegration_F08_ContactToAdminView(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// Visitor kirim pesan
	msg := map[string]string{"Name": "Visitor Budi", "Phone": "08123", "Email": "budi@test.com", "Message": "Saya tertarik"}
	msgJSON, _ := json.Marshal(msg)
	contactReq := httptest.NewRequest(http.MethodPost, "/api/contact", bytes.NewBuffer(msgJSON))
	contactReq.Header.Set("Content-Type", "application/json")
	contactW := httptest.NewRecorder()
	router.ServeHTTP(contactW, contactReq)
	if contactW.Code != http.StatusOK {
		t.Fatalf("Contact expected 200, got %d", contactW.Code)
	}

	// Admin login dan lihat messages
	token := registerAndLogin(t, router, "msgadmin", "msgpass")
	msgReq := httptest.NewRequest(http.MethodGet, "/api/admin/messages", nil)
	msgReq.Header.Set("Authorization", "Bearer "+token)
	msgW := httptest.NewRecorder()
	router.ServeHTTP(msgW, msgReq)

	var messages []models.Message
	json.Unmarshal(msgW.Body.Bytes(), &messages)
	if len(messages) != 1 || messages[0].Name != "Visitor Budi" {
		t.Fatal("Message not visible to admin")
	}
	t.Log("✅ F08-IT01: Alur Contact → Admin View Messages berhasil")
}

// ============================================================
// FITUR 9: Full CRUD Product
// ============================================================

func TestIntegration_F09_FullCRUDProduct(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	token := registerAndLogin(t, router, "crudadmin", "crudpass")

	// CREATE
	req, w := createMultipartRequest(t, http.MethodPost, "/api/admin/products",
		map[string]string{"name": "Granit Hitam", "description": "Premium", "category": "Eksterior"}, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Create expected 200, got %d", w.Code)
	}

	// READ
	getReq := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	var products []models.Product
	json.Unmarshal(getW.Body.Bytes(), &products)
	if len(products) != 1 {
		t.Fatalf("Expected 1 product, got %d", len(products))
	}

	// DELETE
	delReq := httptest.NewRequest(http.MethodDelete, "/api/admin/products/1", nil)
	delReq.Header.Set("Authorization", "Bearer "+token)
	delW := httptest.NewRecorder()
	router.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusOK {
		t.Fatalf("Delete expected 200, got %d", delW.Code)
	}

	t.Log("✅ F09-IT01: Full CRUD Product (Create → Read → Delete) berhasil")
}

// ============================================================
// FITUR 10: Full CRUD Category
// ============================================================

func TestIntegration_F10_FullCRUDCategory(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	token := registerAndLogin(t, router, "catcrud", "catcrud")

	// CREATE
	req, w := createMultipartRequest(t, http.MethodPost, "/api/admin/categories",
		map[string]string{"name": "Marmer", "description": "Batu mewah"}, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Create expected 200, got %d", w.Code)
	}

	// READ
	getReq := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	var cats []models.Category
	json.Unmarshal(getW.Body.Bytes(), &cats)
	if len(cats) != 1 {
		t.Fatalf("Expected 1 category, got %d", len(cats))
	}

	// DELETE
	delReq := httptest.NewRequest(http.MethodDelete, "/api/admin/categories/1", nil)
	delReq.Header.Set("Authorization", "Bearer "+token)
	delW := httptest.NewRecorder()
	router.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusOK {
		t.Fatalf("Delete expected 200, got %d", delW.Code)
	}
	t.Log("✅ F10-IT01: Full CRUD Category (Create → Read → Delete) berhasil")
}

// ============================================================
// FITUR 11: Full CRUD Portfolio
// ============================================================

func TestIntegration_F11_FullCRUDPortfolio(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	token := registerAndLogin(t, router, "portcrud", "portcrud")

	// CREATE
	req, w := createMultipartRequest(t, http.MethodPost, "/api/admin/portfolios",
		map[string]string{"title": "Proyek Villa", "description": "Pemasangan"}, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Create expected 200, got %d", w.Code)
	}

	// READ
	getReq := httptest.NewRequest(http.MethodGet, "/api/portfolios", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	var items []models.Portfolio
	json.Unmarshal(getW.Body.Bytes(), &items)
	if len(items) != 1 {
		t.Fatalf("Expected 1 portfolio, got %d", len(items))
	}

	// DELETE
	delReq := httptest.NewRequest(http.MethodDelete, "/api/admin/portfolios/1", nil)
	delReq.Header.Set("Authorization", "Bearer "+token)
	delW := httptest.NewRecorder()
	router.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusOK {
		t.Fatalf("Delete expected 200, got %d", delW.Code)
	}
	t.Log("✅ F11-IT01: Full CRUD Portfolio (Create → Read → Delete) berhasil")
}

// ============================================================
// FITUR 12: Full CRUD Testimoni
// ============================================================

func TestIntegration_F12_CreateAndDeleteTestimoni(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// CREATE via public API
	req, w := createMultipartRequest(t, http.MethodPost, "/api/testimoni",
		map[string]string{"client_name": "Client1", "testimonial_text": "Sangat bagus!"}, "")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Create testimoni expected 200, got %d", w.Code)
	}

	// READ
	getReq := httptest.NewRequest(http.MethodGet, "/api/testimoni", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	var items []models.Testimoni
	json.Unmarshal(getW.Body.Bytes(), &items)
	if len(items) != 1 {
		t.Fatalf("Expected 1 testimoni, got %d", len(items))
	}

	// DELETE
	delReq := httptest.NewRequest(http.MethodDelete, "/api/testimoni/1", nil)
	delW := httptest.NewRecorder()
	router.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusOK {
		t.Fatalf("Delete expected 200, got %d", delW.Code)
	}
	t.Log("✅ F12-IT01: Full CRUD Testimoni (Create → Read → Delete) berhasil")
}

// ============================================================
// FITUR 13: Get Admin Stats
// ============================================================

func TestIntegration_F13_AdminStats(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	token := registerAndLogin(t, router, "statsadmin", "statspass")

	// Seed beberapa data
	config_db_create(&models.Product{Name: "P1"})
	config_db_create(&models.Message{Name: "M1", Message: "Hi"})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var stats map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &stats)
	if stats["products"].(float64) != 1 {
		t.Fatalf("Expected 1 product in stats")
	}
	t.Log("✅ F13-IT01: Alur Auth → Get Admin Stats berhasil")
}

// ============================================================
// FITUR 14: Get/Delete Messages
// ============================================================

func TestIntegration_F14_ContactThenDeleteMessage(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	// Kirim pesan kontak
	msg := map[string]string{"Name": "Test", "Phone": "08", "Email": "t@t.com", "Message": "Hello"}
	msgJSON, _ := json.Marshal(msg)
	contactReq := httptest.NewRequest(http.MethodPost, "/api/contact", bytes.NewBuffer(msgJSON))
	contactReq.Header.Set("Content-Type", "application/json")
	contactW := httptest.NewRecorder()
	router.ServeHTTP(contactW, contactReq)

	// Admin login, lihat, dan hapus pesan
	token := registerAndLogin(t, router, "deladmin", "delpass")

	getReq := httptest.NewRequest(http.MethodGet, "/api/admin/messages", nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	var messages []models.Message
	json.Unmarshal(getW.Body.Bytes(), &messages)
	if len(messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(messages))
	}

	delReq := httptest.NewRequest(http.MethodDelete, "/api/admin/messages/1", nil)
	delReq.Header.Set("Authorization", "Bearer "+token)
	delW := httptest.NewRecorder()
	router.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusOK {
		t.Fatalf("Delete expected 200, got %d", delW.Code)
	}
	t.Log("✅ F14-IT01: Alur Contact → Admin View → Delete Message berhasil")
}

// ============================================================
// FITUR 15: JWT Middleware Protection
// ============================================================

func TestIntegration_F15_AdminRoutes_NoToken(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	endpoints := []string{"/api/admin/stats", "/api/admin/messages", "/api/admin/visitors"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Endpoint %s: Expected 401 without token, got %d", ep, w.Code)
		}
	}
	t.Log("✅ F15-IT01: Semua admin endpoint dilindungi JWT (tanpa token = 401)")
}

func TestIntegration_F15_AdminRoutes_InvalidToken(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	endpoints := []string{"/api/admin/stats", "/api/admin/messages", "/api/admin/visitors"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Endpoint %s: Expected 401 with invalid token, got %d", ep, w.Code)
		}
	}
	t.Log("✅ F15-IT02: Semua admin endpoint menolak token invalid (401)")
}

func TestIntegration_F15_AdminRoutes_ValidToken(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	token := registerAndLogin(t, router, "jwtadmin", "jwtpass")

	endpoints := []string{"/api/admin/stats", "/api/admin/messages", "/api/admin/visitors"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("Endpoint %s: Expected 200 with valid token, got %d", ep, w.Code)
		}
	}
	t.Log("✅ F15-IT03: Semua admin endpoint dapat diakses dengan token valid (200)")
}

func TestIntegration_F15_CreateProduct_WithoutAuth(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req, w := createMultipartRequest(t, http.MethodPost, "/api/admin/products",
		map[string]string{"name": "Test"}, "")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", w.Code)
	}
	t.Log("✅ F15-IT04: Create product tanpa auth ditolak (401)")
}

// ============================================================
// Helper DB functions
// ============================================================

func configDBFindAdmin(admin *models.Admin, username string) bool {
	return config.DB.Where("username = ?", username).First(admin).Error != nil
}

func config_db_create(value interface{}) {
	config.DB.Create(value)
}
