package models

import (
	"time"
)

type InventoryTransaction struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	ItemID          uint      `gorm:"not null"`
	Item            Item      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransactionType string    `gorm:"type:varchar(20);not null"` 
	Quantity        int       `gorm:"not null"`
	ReferenceType   string    `gorm:"type:varchar(50)"`          
	ReferenceID     uint      
	HandledByID     uint      `gorm:"not null"`                  
	User            User      `gorm:"foreignKey:HandledByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	TransactionDate time.Time `gorm:"autoCreateTime"`
}