package models

import (
	"time"
)

type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	RoleName    string    `gorm:"type:varchar(50);not null;unique"` 
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Users []User `gorm:"foreignKey:RoleID"`
}