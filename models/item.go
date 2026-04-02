package models

import "time"

type Item struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CategoryID   uint      `gorm:"not null"`
	Category     Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ItemCode     string    `gorm:"type:varchar(50);uniqueIndex;not null"` 
	Name         string    `gorm:"type:varchar(150);not null"`
	Description  string    `gorm:"type:text"`
	Unit         string    `gorm:"type:varchar(50);not null"`             
	CurrentStock int       `gorm:"default:0"`                             
	MinimumStock int       `gorm:"default:0"`                            
	QRCodeURL    string    `gorm:"type:varchar(255)"`                  
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}