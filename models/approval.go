package models

import (
	"time"
)

type Approval struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	RequestID  uint      `gorm:"not null"`
	Request    Request   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ApproverID uint      `gorm:"not null"` 
	User       User      `gorm:"foreignKey:ApproverID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Status     string    `gorm:"type:varchar(30);not null"` 
	Comments   string    `gorm:"type:text"`
	ActionDate time.Time `gorm:"autoCreateTime"`
}