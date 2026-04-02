package models

import (
	"time"
)

type Request struct {
	ID            uint            `gorm:"primaryKey;autoIncrement"`
	UserID        uint            `gorm:"not null"` 
	User          User            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	RequestNumber string          `gorm:"type:varchar(50);uniqueIndex;not null"`
	RequestDate   time.Time       `gorm:"autoCreateTime"`
	Status        string          `gorm:"type:varchar(30);default:'Pending_Approval'"` 
	Notes         string          `gorm:"type:text"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime"`

	RequestDetails []RequestDetail `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type RequestDetail struct {
	ID                uint  `gorm:"primaryKey;autoIncrement"`
	RequestID         uint  `gorm:"not null"`
	ItemID            uint  `gorm:"not null"`
	Item              Item  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	QuantityRequested int   `gorm:"not null"`
	QuantityApproved  int   `gorm:"default:0"`
}