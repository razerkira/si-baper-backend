package controllers

import (
	"net/http"
	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetMyProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	if err := config.DB.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{"data": user})
}

type UpdateProfileInput struct {
	FullName   string `json:"full_name" binding:"required"`
	Department string `json:"department" binding:"required"`
	Password   string `json:"password"` 
}

func UpdateMyProfile(c *gin.Context) {
	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	userID, _ := c.Get("user_id")
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	user.FullName = input.FullName
	user.Department = input.Department

	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.PasswordHash = string(hashedPassword) 
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui!"})
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	for i := range users {
		users[i].PasswordHash = "" 
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetRoles(c *gin.Context) {
	var roles []models.Role
	if err := config.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

type UpdateUserInput struct {
	FullName   string `json:"full_name" binding:"required"`
	Department string `json:"department" binding:"required"`
	RoleID     uint   `json:"role_id" binding:"required"`
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id") 
	var input UpdateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	user.FullName = input.FullName
	user.Department = input.Department
	user.RoleID = input.RoleID

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil diperbarui!"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus!"})
}