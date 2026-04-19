package controllers

import (
	"net/http"
	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
)

type SubmissionInput struct {
	ItemName string `json:"item_name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

// CreateSubmission: Menambahkan pengajuan baru
func CreateSubmission(c *gin.Context) {
	var input SubmissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid. Pastikan nama barang diisi dan jumlah lebih dari 0."})
		return
	}

	// Ambil UserID dari token JWT (middleware Auth)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak dapat memverifikasi identitas pengguna"})
		return
	}

	submission := models.Submission{
		UserID:   userID.(uint),
		ItemName: input.ItemName,
		Quantity: input.Quantity,
		Status:   "Menunggu Persetujuan",
	}

	// Simpan ke database
	if err := config.DB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data pengajuan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengajuan barang berhasil dikirim!",
		"data":    submission,
	})
}

// GetMySubmissions: Mengambil riwayat pengajuan milik user yang sedang login
func GetMySubmissions(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var submissions []models.Submission
	// Preload User untuk mendapatkan nama pengaju jika diperlukan nanti
	if err := config.DB.Preload("User").Where("user_id = ?", userID).Order("created_at desc").Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat pengajuan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil riwayat pengajuan",
		"data":    submissions,
	})
}

// GetAllSubmissions: (Khusus Admin/Pimpinan) Mengambil semua pengajuan masuk
func GetAllSubmissions(c *gin.Context) {
	var submissions []models.Submission
	
	if err := config.DB.Preload("User").Order("created_at desc").Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar pengajuan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua daftar pengajuan",
		"data":    submissions,
	})
}