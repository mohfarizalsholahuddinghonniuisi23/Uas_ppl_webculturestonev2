package config

import (
	"culturstone/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := os.Getenv("POSTGRES_URL_NON_POOLING")
	if dsn == "" {
		// Fallback for local development or other environments that might not have the non-pooling URL
		dsn = os.Getenv("POSTGRES_URL")
		if dsn == "" {
			panic("POSTGRES_URL or POSTGRES_URL_NON_POOLING environment variable not set")
		}
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Gagal koneksi database: " + err.Error())
	}

	// AutoMigrate will create/update tables based on the models
	log.Println("Memulai migrasi database... (DINONAKTIFKAN UNTUK DEBUGGING)")
	/*
	err = DB.AutoMigrate(&models.Admin{}, &models.Category{}, &models.ProductCategory{}, &models.Product{}, &models.Portfolio{}, &models.Testimoni{}, &models.Message{}, &models.Visitor{}, &models.User{})
	if err != nil {
		log.Printf("FATAL: Gagal melakukan auto migrate database: %v\n", err)
		panic("Gagal melakukan auto migrate database: " + err.Error())
	}
	*/
	log.Println("Migrasi database berhasil. (DINONAKTIFKAN UNTUK DEBUGGING)")


	// Seeding data is optional and might be better handled by a separate script,
	// but we'll leave it here for now.
	// seedData()
}

func seedData() {
	var c int64
	DB.Model(&models.Product{}).Count(&c)
	if c == 0 {
		log.Println("Melakukan seeding data awal...")
		DB.Create(&models.ProductCategory{Name: "Sanitary"})
		DB.Create(&models.ProductCategory{Name: "Dekorasi"})
		DB.Create(&models.ProductCategory{Name: "Eksterior"})

		cats := []models.Category{
			{Name: "Granit", Description: "Batu keras.", ImageURL: "https://images.unsplash.com/photo-1620302247770-c9a933e2594e?w=500"},
			{Name: "Marmer", Description: "Mewah.", ImageURL: "https://images.unsplash.com/photo-1617791160505-6f00504e35d9?w=500"},
		}
		DB.Create(&cats)

		DB.Create(&models.Product{
			Name: "Wastafel Marmer", Category: "Sanitary",
			ImageURL: "https://images.unsplash.com/photo-1584622050111-993a426fbf0a?w=500&q=80",
			Quality:  "Export", Size: "40cm", Finishing: "Polish",
			Description: "Wastafel batu alam minimalis.",
		})
		log.Println("Seeding data selesai.")
	}
}

