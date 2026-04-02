package controllers

import (
	"net/http"
	"os"
	"time"

	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterInput struct {
	NIP        string `json:"nip" binding:"required"`
	FullName   string `json:"full_name" binding:"required"`
	Department string `json:"department" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	RoleID     uint   `json:"role_id" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	user := models.User{
		NIP:          input.NIP,
		FullName:     input.FullName,
		Department:   input.Department,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		RoleID:       input.RoleID,
		Status:       true, 
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data pengguna. Email atau NIP mungkin sudah terdaftar."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil!",
		"data": gin.H{
			"id":         user.ID,
			"nip":        user.NIP,
			"full_name":  user.FullName,
			"email":      user.Email,
			"department": user.Department,
		},
	})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if !user.Status {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda dinonaktifkan. Silakan hubungi Admin."})
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan token sesi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   tokenString,
		"data": gin.H{
			"id":         user.ID,
			"nip":        user.NIP,
			"full_name":  user.FullName,
			"role_id":    user.RoleID,
		},
	})
}