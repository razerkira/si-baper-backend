package controllers

import (
	"net/http"
	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
)

func GetInventoryLogs(c *gin.Context) {
	var logs []models.InventoryTransaction

	if err := config.DB.Preload("Item").Preload("User").Order("transaction_date desc").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil laporan mutasi barang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil jejak audit mutasi barang",
		"data":    logs,
	})
}