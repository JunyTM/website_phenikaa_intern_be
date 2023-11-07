package model

import (
	"time"

	"gorm.io/gorm"
)

// Bảng thông tin công việc thực tập
type InternJob struct {
	ID             uint   `json:"id" gorm:"autoIncrement"`
	CompanyId      uint   `json:"company_id"`
	Title          string `json:"title"`           // Tên công việc
	JobDesc        string `json:"job_description"` // Mô tả công việc
	Require        string `json:"require"`         // Yêu cầu công việc
	Adress         string `json:"adress"`          // Địa chỉ làm việc
	Benefit        string `json:"benefit"`         // Quyền lợi
	FormOfWork     string `json:"form_of_work"`    // Hình thức làm việc
	FemaleQuantity int    `json:"femaleQuantity"`  // Số lượng nữ
	MaleQuantity   int    `json:"maleQuantity"`    // Số lượng nam
	Salary         uint   `json:"salary"`          // Mức lương
	EndDate        string `json:"end_date"`        // Ngày hết hạn nộp hồ sơ

	CreatedAt time.Time      `json:"createdAt" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updatedAt" swaggerignore:"true"`
}
