package main

import (
	"log"
	"os"
	"si-baper-backend/config"
	"si-baper-backend/routes"
	"si-baper-backend/seeders"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("🚀 [START] Aplikasi dinyalakan...")

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ [INFO] File .env tidak ditemukan (Production Mode).")
	}

	router := gin.Default()

	// Pengaturan CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://si-baper.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rute Root untuk memuaskan Health Check Back4App seketika
	router.GET("/", func(c *gin.Context) {
		c.String(200, "SI-BAPER API is UP and Running!")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong!"})
	})

	// Daftarkan semua rute API utama
	routes.SetupRoutes(router)

	// --- TRIK GOROUTINE: JALANKAN DATABASE DI BACKGROUND ---
	// Agar server bisa langsung membuka Port 8080 dalam 0.1 detik
	go func() {
		log.Println("⏳ [DB] Menghubungkan ke Aiven secara asinkron...")
		config.ConnectDB()
		log.Println("✅ [DB] Terhubung ke Database Aiven!")

		log.Println("⏳ [SEEDER] Menjalankan seeder...")
		seeders.SeedRoles(config.DB)
		seeders.SeedAdminUser(config.DB)
		log.Println("✅ [SEEDER] Seeder Selesai!")
	}()
	// -------------------------------------------------------

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🔥 [SERVER] Langsung membuka port %s...", port)
	router.Run(":" + port)
}