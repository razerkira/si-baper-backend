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
	log.Println("🚀 [START] Aplikasi mulai dinyalakan...")

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ [INFO] File .env tidak ditemukan (wajar di production).")
	}

	log.Println("⏳ [DB] Mencoba menghubungi database Aiven...")
	config.ConnectDB()
	log.Println("✅ [DB] Berhasil terhubung ke database!")

	log.Println("⏳ [SEEDER] Menjalankan seeder...")
	seeders.SeedRoles(config.DB)
	seeders.SeedAdminUser(config.DB)
	log.Println("✅ [SEEDER] Seeder selesai!")

	log.Println("⏳ [ROUTER] Menyiapkan rute API...")
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://si-baper.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong! SI-BAPER API is running smoothly."})
	})

	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🔥 [SERVER] Bersiap mendengarkan di port %s...", port)
	router.Run(":" + port)
}