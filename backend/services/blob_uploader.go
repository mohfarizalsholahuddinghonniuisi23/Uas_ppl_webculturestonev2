package services

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// allowedExtensions adalah whitelist ekstensi file yang diperbolehkan untuk diupload.
// BUG-03 FIX: Hanya file gambar dan video yang diizinkan.
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".mp4":  true,
	".mov":  true,
	".avi":  true,
}

// allowedMimeTypes adalah whitelist MIME type yang diperbolehkan.
var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"video/mp4":       true,
	"video/quicktime": true,
	"video/x-msvideo": true,
}

// UploadFileLocal menyimpan file ke folder ./uploads/ secara lokal
// dan mengembalikan URL akses file tersebut.
// BUG-10 FIX: Nama fungsi diubah dari UploadFileToVercelBlob → UploadFileLocal.
// BUG-03 FIX: Validasi ekstensi & MIME type sebelum menyimpan.
func UploadFileLocal(file *multipart.FileHeader, prefix string) (string, error) {
	// --- BUG-03 FIX: Validasi ekstensi file ---
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("tipe file tidak diizinkan: '%s'. Hanya jpg, jpeg, png, gif, webp, mp4, mov, avi yang diperbolehkan", ext)
	}

	// Buka file yang diupload
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file yang diupload: %w", err)
	}
	defer src.Close()

	// --- BUG-03 FIX: Validasi MIME type dari isi file (bukan hanya ekstensi) ---
	// Baca 512 byte pertama untuk deteksi MIME type
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("gagal membaca file untuk validasi MIME: %w", err)
	}
	detectedMime := http.DetectContentType(buffer[:n])
	// Ambil hanya bagian utama MIME (sebelum ";")
	mimeBase := strings.Split(detectedMime, ";")[0]
	mimeBase = strings.TrimSpace(mimeBase)
	if !allowedMimeTypes[mimeBase] {
		return "", fmt.Errorf("konten file tidak diizinkan (MIME: %s). Hanya file gambar dan video yang diperbolehkan", mimeBase)
	}

	// Reset reader agar file bisa dibaca ulang dari awal
	if _, err = src.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("gagal me-reset posisi file: %w", err)
	}

	// Buat folder uploads jika belum ada
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("gagal membuat folder uploads: %w", err)
	}

	// Buat nama file yang unik agar tidak overwrite
	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	destPath := filepath.Join(uploadsDir, filename)

	// Buat file tujuan
	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file tujuan '%s': %w", destPath, err)
	}
	defer dst.Close()

	// Salin isi file
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("gagal menyalin file: %w", err)
	}

	// BUG-10 FIX: Ambil base URL dari environment variable agar bisa dikonfigurasi.
	// Fallback ke localhost:8080 jika BASE_URL tidak di-set.
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	fileURL := fmt.Sprintf("%s/uploads/%s", baseURL, filename)
	log.Printf("File '%s' berhasil disimpan lokal. URL: %s", filename, fileURL)

	return fileURL, nil
}

// UploadFileToVercelBlob adalah alias untuk UploadFileLocal untuk backward-compatibility.
// Deprecated: Gunakan UploadFileLocal langsung.
func UploadFileToVercelBlob(file *multipart.FileHeader, prefix string) (string, error) {
	return UploadFileLocal(file, prefix)
}