package controllers

import (
	"net/http"

	"si-baper-backend/config"
	"si-baper-backend/models"

	"github.com/gin-gonic/gin"
)

type ApprovalDetailInput struct {
	RequestDetailID  uint `json:"request_detail_id" binding:"required"`
	QuantityApproved int  `json:"quantity_approved" binding:"min=0"`
}

type ProcessApprovalInput struct {
	RequestID uint                  `json:"request_id" binding:"required"`
	Status    string                `json:"status" binding:"required,oneof=Approved Rejected"`
	Comments  string                `json:"comments"`
	Items     []ApprovalDetailInput `json:"items"` 
}

func GetPendingRequests(c *gin.Context) {
	var requests []models.Request

	if err := config.DB.Where("status = ?", "Pending_Approval").
		Preload("User"). 
		Preload("RequestDetails.Item"). 
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data antrean persetujuan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil antrean persetujuan",
		"data":    requests,
	})
}

func ProcessApproval(c *gin.Context) {
	var input ProcessApprovalInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	approverID, _ := c.Get("user_id")

	tx := config.DB.Begin()

	var request models.Request
	if err := tx.Where("id = ? AND status = ?", input.RequestID, "Pending_Approval").First(&request).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Formulir permintaan tidak ditemukan atau sudah diproses"})
		return
	}

	approval := models.Approval{
		RequestID:  request.ID,
		ApproverID: uint(approverID.(float64)),
		Status:     input.Status,
		Comments:   input.Comments,
	}
	if err := tx.Create(&approval).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan riwayat persetujuan"})
		return
	}

	if input.Status == "Approved" {
		for _, itemInput := range input.Items {
			var detail models.RequestDetail
			
			if err := tx.Preload("Item").First(&detail, itemInput.RequestDetailID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Detail barang tidak ditemukan dalam formulir ini"})
				return
			}

			if itemInput.QuantityApproved > detail.Item.CurrentStock {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak mencukupi untuk barang: " + detail.Item.Name})
				return
			}

			detail.QuantityApproved = itemInput.QuantityApproved
			if err := tx.Save(&detail).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui jumlah barang yang disetujui"})
				return
			}

			detail.Item.CurrentStock -= itemInput.QuantityApproved
			if err := tx.Save(&detail.Item).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memotong stok barang otomatis"})
				return
			}

			inventoryLog := models.InventoryTransaction{
				ItemID:          detail.ItemID,
				TransactionType: "OUT",                
				Quantity:        itemInput.QuantityApproved,
				ReferenceType:   "Approval",
				ReferenceID:     request.ID,
				HandledByID:     uint(approverID.(float64)),
			}
			if err := tx.Create(&inventoryLog).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat log mutasi barang"})
				return
			}
		}
	}

	request.Status = input.Status
	if err := tx.Save(&request).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status permintaan"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Persetujuan berhasil diproses!",
		"data": gin.H{
			"request_number": request.RequestNumber,
			"new_status":     request.Status,
		},
	})
}