package controllers

import (
	"net/http"
	"si-baper-backend/config"
	"si-baper-backend/models"
	"si-baper-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode" // <-- Import library pembuat QR Code
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

	// 2. Buat gambar QR Code (mengubah text ItemCode menjadi data gambar byte array)
	pngData, err := qrcode.Encode(item.ItemCode, qrcode.Medium, 256)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan gambar QR Code: " + err.Error()})
		return
	}

	// 3. Unggah data gambar tersebut ke Cloudinary
	qrCodeURL, err := utils.UploadQRCodeToCloudinary(item.ItemCode, pngData)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah QR Code ke Cloudinary: " + err.Error()})
		return
	}

	// 4. Simpan URL aman dari Cloudinary kembali ke database barang
	item.QRCodeURL = qrCodeURL
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan URL QR Code ke database"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Barang dan QR Code berhasil ditambahkan!",
		"data":    item,
	})
}