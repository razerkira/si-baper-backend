package controllers

import (
	"net/http"
	"si-baper-backend/config"
	"si-baper-backend/models"
	"si-baper-backend/utils"

	"github.com/gin-gonic/gin"
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

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback() 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan barang. Kode Barang unik atau Kategori valid."})
		return
	}

	qrCodeURL, err := utils.GenerateAndUploadQRCode(item.ItemCode)
	if err != nil {
		tx.Rollback() 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan QR Code: " + err.Error()})
		return
	}

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