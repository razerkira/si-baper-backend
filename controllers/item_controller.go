// src/controllers/item_controller.go
package controllers

import (
	"fmt" // <-- Tambahkan import fmt untuk memformat teks string
	"net/http"
	"os"
	"si-baper-backend/config"
	"si-baper-backend/models"
	"si-baper-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type ItemInput struct {
	CategoryID   uint   `json:"category_id" binding:"required"`
	ItemCode     string `json:"item_code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Unit         string `json:"unit" binding:"required"`
	CurrentStock int    `json:"current_stock"`
	MinimumStock int    `json:"minimum_stock"`
}

func GetItems(c *gin.Context) {
	var items []models.Item

	if err := config.DB.Preload("Category").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data barang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data katalog",
		"data":    items,
	})
}

func CreateItem(c *gin.Context) {
	var input ItemInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	item := models.Item{
		CategoryID:   input.CategoryID,
		ItemCode:     input.ItemCode,
		Name:         input.Name,
		Description:  input.Description,
		Unit:         input.Unit,
		CurrentStock: input.CurrentStock,
		MinimumStock: input.MinimumStock,
	}

	// 1. Simpan data awal barang ke database
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan barang. Kode Barang mungkin sudah ada."})
		return
	}

	// 2. Siapkan konten lengkap untuk QR Code
	qrContent := fmt.Sprintf(
		"Kode Barang: %s\nNama Barang: %s\nKategori ID: %d\nStok Saat Ini: %d %s\nBatas Minimum: %d %s\nDeskripsi: %s",
		item.ItemCode,
		item.Name,
		item.CategoryID,
		item.CurrentStock, item.Unit,
		item.MinimumStock, item.Unit,
		item.Description,
	)

	// 3. Buat gambar QR Code menggunakan konten lengkap
	pngData, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan gambar QR Code: " + err.Error()})
		return
	}

	// 4. Unggah data gambar tersebut ke Local Storage VPS
	qrCodePath, err := utils.SaveQRCodeLocally(item.ItemCode, pngData)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file QR Code di server lokal: " + err.Error()})
		return
	}

	// 5. Simpan path URL ke database barang
	item.QRCodeURL = qrCodePath
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan path URL QR Code ke database"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Barang dan QR Code berhasil ditambahkan!",
		"data":    item,
	})
}

func UpdateItem(c *gin.Context) {
	itemID := c.Param("id")
	var item models.Item

	// 1. Cari barang di database
	if err := config.DB.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	var input ItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	// 2. Hapus file QR Code lama (karena datanya pasti berubah, kita hapus dulu yang lama)
	if item.QRCodeURL != "" {
		_ = os.Remove(item.QRCodeURL) // Mengabaikan error jika file tidak ditemukan
	}

	// 3. Siapkan konten lengkap yang BARU untuk QR Code
	newQrContent := fmt.Sprintf(
		"Kode Barang: %s\nNama Barang: %s\nKategori ID: %d\nStok Saat Ini: %d %s\nBatas Minimum: %d %s\nDeskripsi: %s",
		input.ItemCode,
		input.Name,
		input.CategoryID,
		input.CurrentStock, input.Unit,
		input.MinimumStock, input.Unit,
		input.Description,
	)

	// 4. Buat gambar QR Code baru
	pngData, err := qrcode.Encode(newQrContent, qrcode.Medium, 256)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan gambar QR Code baru: " + err.Error()})
		return
	}

	// 5. Simpan file QR Code baru ke Local Storage VPS
	qrCodePath, err := utils.SaveQRCodeLocally(input.ItemCode, pngData)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file QR Code baru: " + err.Error()})
		return
	}

	// 6. Perbarui semua field barang
	item.QRCodeURL = qrCodePath
	item.CategoryID = input.CategoryID
	item.ItemCode = input.ItemCode
	item.Name = input.Name
	item.Description = input.Description
	item.Unit = input.Unit
	item.CurrentStock = input.CurrentStock
	item.MinimumStock = input.MinimumStock

	// 7. Simpan perubahan ke database
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data barang. Pastikan Kode Barang unik."})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Data barang berhasil diperbarui!",
		"data":    item,
	})
}

func DeleteItem(c *gin.Context) {
	itemID := c.Param("id")
	var item models.Item

	if err := config.DB.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus barang dari database"})
		return
	}

	if item.QRCodeURL != "" {
		err := os.Remove(item.QRCodeURL)
		if err != nil {
			// fmt.Println("Gagal menghapus file fisik QR Code:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Barang dan file QR Code berhasil dihapus!",
	})
}