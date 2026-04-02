package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	RoleID       uint           `gorm:"not null"`
	Role         Role           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	NIP          string         `gorm:"type:varchar(50);uniqueIndex;not null"` 
	FullName     string         `gorm:"type:varchar(150);not null"`
	Department   string         `gorm:"type:varchar(150);not null"`            
	Email        string         `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`            
	Status       bool           `gorm:"default:true"`                          
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`                                 
}