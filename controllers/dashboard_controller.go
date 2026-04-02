package controllers

import (
	"net/http"
	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	var totalItems int64
	var pendingRequests int64
	var approvedRequests int64
	var totalRequests int64

	config.DB.Model(&models.Item{}).Count(&totalItems)

	config.DB.Model(&models.Request{}).Where("status = ?", "Pending_Approval").Count(&pendingRequests)
	config.DB.Model(&models.Request{}).Where("status = ?", "Approved").Count(&approvedRequests)
	
	config.DB.Model(&models.Request{}).Count(&totalRequests)

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil statistik dashboard",
		"data": gin.H{
			"total_items":       totalItems,
			"pending_requests":  pendingRequests,
			"approved_requests": approvedRequests,
			"total_requests":    totalRequests,
		},
	})
}