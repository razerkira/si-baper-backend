package models

import "time"

// Submission merepresentasikan pengajuan pengadaan barang
type Submission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ItemName  string    `gorm:"not null" json:"item_name"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Status    string    `gorm:"default:'Menunggu Persetujuan'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}