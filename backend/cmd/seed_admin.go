// Jalankan perintah ini dari folder backend/:
//   go run cmd/seed_admin.go
//
// Script ini membuat akun admin pertama di database.
// Jika sudah ada admin, script akan menolak dan memberitahu.

package main

import (
	"bufio"
	"culturstone/config"
	"culturstone/models"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("  CULTURESTONE — Setup Admin Pertama   ")
	fmt.Println("========================================")

	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("[INFO] File .env tidak ditemukan, menggunakan environment variable sistem.")
	}

	// 2. Konek ke database
	fmt.Println("\n[1/4] Menghubungkan ke database...")
	config.ConnectDB()
	fmt.Println("      ✓ Berhasil konek ke database.")

	// 3. Cek apakah sudah ada admin
	fmt.Println("[2/4] Memeriksa data admin yang sudah ada...")
	var count int64
	config.DB.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		fmt.Println("\n❌ GAGAL: Sudah ada admin di database.")
		fmt.Println("   Gunakan akun admin yang sudah ada untuk login.")
		fmt.Println("   Jika lupa password, hubungi developer untuk reset manual.")
		os.Exit(1)
	}
	fmt.Println("      ✓ Belum ada admin. Lanjut pendaftaran.")

	// 4. Input username
	fmt.Println("[3/4] Masukkan data admin baru:")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("      Username : ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "" {
		fmt.Println("\n❌ GAGAL: Username tidak boleh kosong.")
		os.Exit(1)
	}

	// 5. Input password
	// Catatan: Password akan terlihat di terminal saat diketik.
	// Ini aman karena hanya dijalankan oleh developer secara lokal.
	fmt.Print("      Password : ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// 6. Validasi kekuatan password
	fmt.Println("[4/4] Memvalidasi dan menyimpan admin...")
	if err := validatePassword(password); err != nil {
		fmt.Println("\n❌ GAGAL:", err.Error())
		fmt.Println("\nKetentuan password:")
		fmt.Println("  - Minimal 8 karakter")
		fmt.Println("  - Mengandung huruf kapital (A-Z)")
		fmt.Println("  - Mengandung huruf kecil (a-z)")
		fmt.Println("  - Mengandung angka (0-9)")
		fmt.Println("\nContoh password yang valid: Admin2024")
		os.Exit(1)
	}

	// 7. Hash password dengan bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("\n❌ GAGAL hash password:", err.Error())
		os.Exit(1)
	}

	// 8. Simpan ke database
	admin := models.Admin{
		Username: username,
		Password: string(hashed),
	}
	if err := config.DB.Create(&admin).Error; err != nil {
		fmt.Println("\n❌ GAGAL menyimpan admin ke database:", err.Error())
		os.Exit(1)
	}

	fmt.Println("\n========================================")
	fmt.Printf("  ✅ BERHASIL! Admin '%s' sudah terdaftar.\n", username)
	fmt.Println("  Sekarang login di: http://localhost:5173/login")
	fmt.Println("========================================")
}

// validatePassword memvalidasi kekuatan password.
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
		return fmt.Errorf("password harus mengandung minimal 1 huruf kapital (A-Z)")
	}
	if !hasLower {
		return fmt.Errorf("password harus mengandung minimal 1 huruf kecil (a-z)")
	}
	if !hasDigit {
		return fmt.Errorf("password harus mengandung minimal 1 angka (0-9)")
	}
	return nil
}
