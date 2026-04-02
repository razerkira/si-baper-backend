package seeders

import (
	"log"
	"si-baper-backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@kemenham.go.id").Count(&count)

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("4dm1nBmN1TJ3N@2026"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Gagal melakukan hashing password: %v", err)
		}

		adminUser := models.User{
			NIP:          "1111111111",
			FullName:     "admin",
			Email:        "admin@kemenham.go.id",
			Department:   "BMN",
			PasswordHash: string(hashedPassword),
			RoleID:       1,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatalf("Gagal menjalankan seeder admin: %v", err)
		}

		log.Println("Seeder: User Admin BMN berhasil ditambahkan!")
	} else {
		log.Println("Seeder: User Admin BMN sudah ada di database. Dilewati.")
	}
}