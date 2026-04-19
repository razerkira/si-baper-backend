package controllers

import (
	"fmt"
	"net/http"
	"time"

	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
)

// --- STRUKTUR INPUT STANDAR ---
type RequestDetailInput struct {
	ItemID            uint `json:"item_id" binding:"required"`
	QuantityRequested int  `json:"quantity_requested" binding:"required,min=1"`
}

type CreateRequestInput struct {
	Notes string               `json:"notes"`
	Items []RequestDetailInput `json:"items" binding:"required,min=1"` 
}

// --- STRUKTUR INPUT KHUSUS PUBLIK (DENGAN NIP) ---
type PublicRequestInput struct {
	NIP   string               `json:"nip" binding:"required"`
	Notes string               `json:"notes"`
	Items []RequestDetailInput `json:"items" binding:"required,min=1"` 
}

// ------------------------------------------------------------------
// 1. CreateRequest: Untuk Pegawai yang sudah Log In
// ------------------------------------------------------------------
func CreateRequest(c *gin.Context) {
	var input CreateRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid atau barang kosong"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	dateStr := time.Now().Format("20060102")
	randomNum := time.Now().UnixNano() % 100000
	requestNumber := fmt.Sprintf("REQ-%s-%05d", dateStr, randomNum)

	tx := config.DB.Begin()

	request := models.Request{
		UserID:        uint(userID.(float64)), 
		RequestNumber: requestNumber,
		Status:        "Pending_Approval",
		Notes:         input.Notes,
	}

	if err := tx.Create(&request).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat formulir permintaan"})
		return
	}

	for _, itemInput := range input.Items {
		detail := models.RequestDetail{
			RequestID:         request.ID,
			ItemID:            itemInput.ItemID,
			QuantityRequested: itemInput.QuantityRequested,
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan detail barang permintaan"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Permintaan berhasil diajukan dan menunggu persetujuan!",
		"data": gin.H{
			"request_number": request.RequestNumber,
			"request_id":     request.ID,
			"status":         request.Status,
		},
	})
}

// ------------------------------------------------------------------
// 2. GetMyRequests: Untuk melihat riwayat permintaan user yang Log In
// ------------------------------------------------------------------
func GetMyRequests(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var requests []models.Request

	if err := config.DB.Where("user_id = ?", userID).Preload("RequestDetails.Item").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat permintaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil riwayat permintaan",
		"data":    requests,
	})
}

// ------------------------------------------------------------------
// 3. CreatePublicRequest: Untuk Pegawai dari Beranda (Tanpa Log In)
// ------------------------------------------------------------------
func CreatePublicRequest(c *gin.Context) {
	var input PublicRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid. Pastikan NIP dan Barang diisi."})
		return
	}

	// CARI USER BERDASARKAN NIP
	var user models.User
	if err := config.DB.Where("nip = ?", input.NIP).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NIP tidak terdaftar di sistem kami."})
		return
	}

	// Generate Nomor Tiket
	dateStr := time.Now().Format("20060102")
	randomNum := time.Now().UnixNano() % 100000
	requestNumber := fmt.Sprintf("REQ-%s-%05d", dateStr, randomNum)

	tx := config.DB.Begin()

	// Gunakan UserID dari hasil pencarian NIP
	request := models.Request{
		UserID:        user.ID, 
		RequestNumber: requestNumber,
		Status:        "Pending_Approval",
		Notes:         input.Notes,
	}

	if err := tx.Create(&request).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat formulir permintaan publik"})
		return
	}

	for _, itemInput := range input.Items {
		detail := models.RequestDetail{
			RequestID:         request.ID,
			ItemID:            itemInput.ItemID,
			QuantityRequested: itemInput.QuantityRequested,
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan detail barang permintaan"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Permintaan berhasil diajukan dari Beranda!",
		"data": gin.H{
			"request_number": request.RequestNumber,
			"request_id":     request.ID,
			"status":         request.Status,
		},
	})
}