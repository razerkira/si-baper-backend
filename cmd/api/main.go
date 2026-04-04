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
	log.Println("========================================")
	log.Println("🚀 [START] Aplikasi SI-BAPER mulai dinyalakan...")
	log.Println("========================================")

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ [INFO] File .env tidak ditemukan (Production Mode / Environment System).")
	}

	router := gin.Default()

	// --- PENGATURAN CORS (JURUS SAPU JAGAT) ---
	// Menerima request dari domain apapun agar Vercel bisa terhubung dengan mulus
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // Izinkan semua origin
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// ------------------------------------------

	// --- RUTE STATIS UNTUK GAMBAR (SANGAT PENTING) ---
	// Mengizinkan frontend mengakses folder uploads untuk melihat gambar QR Code
	// Contoh akses: http://IP_VPS:8080/uploads/qrcodes/ITEM001.png
	router.Static("/uploads", "./uploads")
	// -------------------------------------------------

	// Rute Root & Ping
	router.GET("/", func(c *gin.Context) {
		c.String(200, "SI-BAPER API is UP and Running on VPS!")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong! Server SI-BAPER Sehat."})
	})

	// Daftarkan semua rute API utama
	routes.SetupRoutes(router)

	// --- TRIK GOROUTINE DENGAN TAMENG ANTI-CRASH ---
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("🚨 [CRASH PREVENTED] Terjadi kegagalan saat inisialisasi:", r)
			}
		}()

		log.Println("⏳ [DB] Menghubungkan ke Database SI-BAPER (Local)...")
		config.ConnectDB()
		log.Println("✅ [DB] Berhasil terhubung ke Database SI-BAPER!")

		log.Println("⏳ [SEEDER] Menjalankan seeder...")
		seeders.SeedRoles(config.DB)
		seeders.SeedAdminUser(config.DB)
		log.Println("✅ [SEEDER] Seeder Selesai!")
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// BIND EKSPLISIT KE 0.0.0.0 AGAR BISA DIAKSES DARI INTERNET
	log.Printf("🔥 [SERVER] Mendengarkan di port %s...", port)
	router.Run("0.0.0.0:" + port)
}