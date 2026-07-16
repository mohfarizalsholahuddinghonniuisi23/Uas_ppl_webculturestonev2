package config

import (
	"culturstone/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	// Baca konfigurasi dari environment variable
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "culturestone"
	}

	// Format DSN untuk MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Gagal koneksi database MySQL: " + err.Error())
	}

	log.Printf("Berhasil konek ke MySQL: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	// AutoMigrate: buat/update tabel otomatis sesuai model
	log.Println("Memulai AutoMigrate...")
	err = DB.AutoMigrate(
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
		log.Printf("PERINGATAN: AutoMigrate gagal: %v\n", err)
	} else {
		log.Println("AutoMigrate selesai — semua tabel siap.")
	}
}
