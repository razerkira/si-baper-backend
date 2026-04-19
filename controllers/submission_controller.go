package controllers

import (
	"net/http"
	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
)

// --- STRUKTUR INPUT STANDAR ---
type SubmissionInput struct {
	ItemName string `json:"item_name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

// --- STRUKTUR INPUT KHUSUS PUBLIK (DENGAN NIP) ---
type PublicSubmissionInput struct {
	NIP      string `json:"nip" binding:"required"`
	ItemName string `json:"item_name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

// ------------------------------------------------------------------
// 1. CreateSubmission: Menambahkan pengajuan baru (Dengan Log In)
// ------------------------------------------------------------------
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

// ------------------------------------------------------------------
// 2. GetMySubmissions: Mengambil riwayat pengajuan milik user yang sedang login
// ------------------------------------------------------------------
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

// ------------------------------------------------------------------
// 3. GetAllSubmissions: (Khusus Admin/Pimpinan) Mengambil semua pengajuan masuk
// ------------------------------------------------------------------
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

// ------------------------------------------------------------------
// 4. CreatePublicSubmission: Untuk Pegawai dari Beranda (Tanpa Log In)
// ------------------------------------------------------------------
func CreatePublicSubmission(c *gin.Context) {
	var input PublicSubmissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid. Pastikan NIP, nama barang, dan jumlah diisi dengan benar."})
		return
	}

	// CARI USER BERDASARKAN NIP
	var user models.User
	if err := config.DB.Where("nip = ?", input.NIP).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NIP tidak terdaftar di sistem kami."})
		return
	}

	submission := models.Submission{
		UserID:   user.ID, // Hubungkan dengan ID dari NIP yang ditemukan
		ItemName: input.ItemName,
		Quantity: input.Quantity,
		Status:   "Menunggu Persetujuan",
	}

	if err := config.DB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengajuan publik"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengajuan berhasil dikirim dari Beranda!",
		"data":    submission,
	})
}