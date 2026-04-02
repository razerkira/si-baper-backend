package controllers

import (
	"net/http"
	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
)

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func GetCategories(c *gin.Context) {
	var categories []models.Category

	if err := config.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data kategori",
		"data":    categories,
	})
}

func CreateCategory(c *gin.Context) {
	var input CategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{Name: input.Name}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori. Nama mungkin sudah ada."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori berhasil ditambahkan!",
		"data":    category,
	})
}