package seeders

import (
	"fmt"
	"si-baper-backend/models"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{RoleName: "Admin", Description: "Administrator Sistem & Gudang"},
		{RoleName: "Pegawai", Description: "Staf pemohon barang persediaan"},
		{RoleName: "Eksekutif", Description: "Pimpinan (hanya melihat dashboard)"},
	}

	for _, role := range roles {
		db.FirstOrCreate(&role, models.Role{RoleName: role.RoleName})
	}

	fmt.Println("Seeder: Data Roles berhasil disiapkan!")
}
