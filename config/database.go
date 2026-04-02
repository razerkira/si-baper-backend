package config

import (
	"fmt"
	"log"
	"os"
	"si-baper-backend/models" 

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT") 

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database!\n", err)
	}

	fmt.Println("Koneksi database PostgreSQL berhasil di port", port)

	err = database.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Category{}, 
		&models.Item{},     
		&models.Request{},       
		&models.RequestDetail{}, 
		&models.Approval{},
		&models.InventoryTransaction{},
	)

	if err != nil {
		log.Fatal("Gagal melakukan Auto-Migrate: ", err)
	}

	fmt.Println("Auto-Migrate berhasil! Tabel roles dan users sudah siap digunakan.")

	DB = database
}