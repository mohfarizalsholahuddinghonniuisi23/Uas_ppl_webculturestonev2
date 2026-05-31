package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ============================================================
// Unit Test: GenerateToken
// ============================================================

// UT-T01: Generate token untuk user_id = 1
func TestGenerateToken_ValidUserID(t *testing.T) {
	token, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token string")
	}
	t.Logf("✅ UT-T01 PASS: Token generated successfully for user_id=1, token length=%d", len(token))
}

// UT-T02: Generate token untuk user_id = 0
func TestGenerateToken_ZeroUserID(t *testing.T) {
	token, err := GenerateToken(0)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token string even for user_id=0")
	}
	t.Logf("✅ UT-T02 PASS: Token generated for user_id=0, token length=%d", len(token))
}

// UT-T03: Generate token untuk user_id besar (999999)
func TestGenerateToken_LargeUserID(t *testing.T) {
	token, err := GenerateToken(999999)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token string for large user_id")
	}
	t.Logf("✅ UT-T03 PASS: Token generated for user_id=999999, token length=%d", len(token))
}

// ============================================================
// Unit Test: ExtractToken
// ============================================================

// UT-T04: Header "Bearer validtoken123"
func TestExtractToken_ValidBearerToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer validtoken123")

	result := ExtractToken(c)
	if result != "validtoken123" {
		t.Fatalf("Expected 'validtoken123', got: '%s'", result)
	}
	t.Logf("✅ UT-T04 PASS: ExtractToken correctly returned 'validtoken123'")
}

// UT-T05: Header kosong
func TestExtractToken_EmptyHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	// No Authorization header set

	result := ExtractToken(c)
	if result != "" {
		t.Fatalf("Expected empty string, got: '%s'", result)
	}
	t.Logf("✅ UT-T05 PASS: ExtractToken returned empty string for empty header")
}

// UT-T06: Header "InvalidFormat" (tanpa "Bearer " prefix)
func TestExtractToken_InvalidFormat(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	result := ExtractToken(c)
	if result != "" {
		t.Fatalf("Expected empty string for invalid format, got: '%s'", result)
	}
	t.Logf("✅ UT-T06 PASS: ExtractToken returned empty string for 'InvalidFormat'")
}

// UT-T07: Header "Bearer" saja (tanpa token setelahnya)
func TestExtractToken_BearerWithoutToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer")

	result := ExtractToken(c)
	if result != "" {
		t.Fatalf("Expected empty string for 'Bearer' only, got: '%s'", result)
	}
	t.Logf("✅ UT-T07 PASS: ExtractToken returned empty string for 'Bearer' only")
}

// ============================================================
// Unit Test: TokenValid
// ============================================================

// UT-T08: Request dengan token valid (dari GenerateToken)
func TestTokenValid_WithValidToken(t *testing.T) {
	// Generate a valid token first
	token, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	err = TokenValid(c)
	if err != nil {
		t.Fatalf("Expected token to be valid, got error: %v", err)
	}
	t.Logf("✅ UT-T08 PASS: TokenValid correctly validated a valid token")
}

// UT-T09: Request tanpa header Authorization
func TestTokenValid_WithoutAuthHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	err := TokenValid(c)
	if err == nil {
		t.Fatal("Expected error for missing Authorization header, got nil")
	}
	t.Logf("✅ UT-T09 PASS: TokenValid returned error for missing auth: %v", err)
}

// UT-T10: Request dengan token invalid/rusak
func TestTokenValid_WithInvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid.token.here")

	err := TokenValid(c)
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
	t.Logf("✅ UT-T10 PASS: TokenValid returned error for invalid token: %v", err)
}
