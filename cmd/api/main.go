package main

import (
	"log"
	"os"
	"time" // <-- Package time untuk jebakan waktu

	"si-baper-backend/config"
	"si-baper-backend/routes"
	"si-baper-backend/seeders"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// --- JEBAKAN WAKTU ---
	// Tahan aplikasi selama 3 detik agar sistem log Back4App sempat menyala dan merekam!
	log.Println("========================================")
	log.Println("🚀 [START] Aplikasi mulai dinyalakan...")
	log.Println("========================================")
	time.Sleep(3 * time.Second)

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ [INFO] File .env tidak ditemukan (Production Mode).")
	}

	router := gin.Default()

	// --- PENGATURAN CORS (JURUS SAPU JAGAT) ---
	// Menerima request dari domain apapun untuk mengatasi error CORS
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // Izinkan semua origin sementara waktu
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// ------------------------------------------

	// Rute Root untuk memuaskan Health Check Back4App
	router.GET("/", func(c *gin.Context) {
		c.String(200, "SI-BAPER API is UP and Running!")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong!"})
	})

	routes.SetupRoutes(router)

	// --- TRIK GOROUTINE DENGAN TAMENG ANTI-CRASH ---
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("🚨 [CRASH PREVENTED] Terjadi kegagalan saat inisialisasi:", r)
			}
		}()

		log.Println("⏳ [DB] Menghubungkan ke Aiven secara asinkron...")
		config.ConnectDB()
		log.Println("✅ [DB] Terhubung ke Database Aiven!")

		log.Println("⏳ [SEEDER] Menjalankan seeder...")
		seeders.SeedRoles(config.DB)
		seeders.SeedAdminUser(config.DB)
		log.Println("✅ [SEEDER] Seeder Selesai!")
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// BIND EKSPLISIT KE 0.0.0.0 AGAR DOCKER BISA MENGAKSESNYA DARI LUAR
	log.Printf("🔥 [SERVER] Langsung membuka port %s...", port)
	router.Run("0.0.0.0:" + port)
}