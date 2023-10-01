package model

import (
	"time"

	"gorm.io/gorm"
)

type InternShipEvaluate struct {
	ID         uint   `json:"id" gorm:"autoIncrement"`
	Purpose    string `json:"purpose"`    // Mục đích thực tập
	Attitude   string `json:"attitude"`   // Thái độ
	Ability    string `json:"ability"`    // Năng lực
	Knowledge  string `json:"knowledge"`  // Kiến thức
	Position   string `json:"position"`   // Vị trí
	JobDesc    string `json:"job_desc"`   // Mô tả công việc
	Instructor string `json:"instructor"` // Người hướng dẫn
	Review     string `json:"review"`     // Nhận xét
	Evaluation string `json:"evaluation"` // Đánh giá
	Note       string `json:"note"`       // Ghi chú
	StartDate  string `json:"start_date"` // Ngày bắt đầu
	EndDate    string `json:"end_date"`   // Ngày kết thúc

	CreatedAt time.Time      `json:"createdAt" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updatedAt" swaggerignore:"true"`
}
