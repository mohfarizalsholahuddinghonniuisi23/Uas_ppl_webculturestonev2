package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Gunakan environment variable untuk secret key di produksi
// Jika tidak ada, gunakan fallback default (HANYA UNTUK DEVELOPMENT)
var apiSecret = []byte(getSecretKey())

func getSecretKey() string {
	secret := os.Getenv("API_SECRET")
	if secret == "" {
		// Fallback ke default yang konsisten
		secret = "RAHASIA_SUPER_AMAN_CULTURSTONE"
	}
	return secret
}

// GenerateToken membuat token JWT baru untuk user tertentu
func GenerateToken(user_id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token berlaku 24 jam

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(apiSecret)
}

// TokenValid memeriksa apakah token dalam request valid
func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	if tokenString == "" {
		return fmt.Errorf("no token provided in Authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan method signing adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return apiSecret, nil
	})

	if err != nil {
		return fmt.Errorf("token parsing error: %v", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// ExtractToken mengambil token string dari header Authorization
// Format yang diharapkan: "Bearer <token>"
func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	// Cek apakah formatnya "Bearer <token>"
	if len(strings.Split(bearerToken, " ")) == 2 {
		token := strings.Split(bearerToken, " ")[1]
		return token
	}
	return ""
}
