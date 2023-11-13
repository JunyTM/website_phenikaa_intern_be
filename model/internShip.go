package model

import (
	"time"

	"gorm.io/gorm"
)

// Bảng quản lý thông tin báo cáo thực tập của sinh viên sinh
type InternShip struct {
	ID                   uint   `json:"id" gorm:"primaryKey"`
	ProfileId            uint   `json:"profile_id"`
	CompanyId            uint   `json:"company_id"`
	InternshipEvaluateId uint   `json:"internship_evaluate_id"`
	Code                 string `json:"code"` // Trang thái đã hoàn thành hay chưa

	// Profile            *Profile            `json:"profile" gorm:"foreignKey:user_id"`
	// Company            *Company            `json:"company" gorm:"foreignKey:CompanyId"`
	// InternShipEvaluate *InternshipEvaluate `json:"internship_evaluate" gorm:"foreignKey:InternshipEvaluateId"`
	CreatedAt          time.Time           `json:"createdAt" swaggerignore:"true"`
	DeletedAt          gorm.DeletedAt      `json:"-" swaggerignore:"true"`
	UpdatedAt          time.Time           `json:"updatedAt" swaggerignore:"true"`
}
