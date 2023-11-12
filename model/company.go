package model

import (
	"time"

	"gorm.io/gorm"
)

// Bảng quản lý thông tin công ty
type Company struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FoundingDay string `json:"founding_day"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Adress      string `json:"adress"`

	CreatedAt time.Time      `json:"createdAt" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updatedAt" swaggerignore:"true"`
}
