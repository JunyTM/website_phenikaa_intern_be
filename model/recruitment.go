package model

import (
	"time"

	"gorm.io/gorm"
)

// Bảng quản lý thông tin ung tuyển thực tập
type Recruitment struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	ProfileId   uint       `json:"profile_id"`    // Id của sinh viên
	InternJobId uint       `json:"intern_job_id"` // Id của bài đăng tuyển dụng
	ProfilePath string     `json:"profile_path"`  // Đường dẫn file CV
	Accepted    bool       `json:"accepted"`      // Đã được nhận hay chưa
	State       string     `json:"state"`         // Trạng thái đã hoàn thành hay chưa
	Profile     *Profile   `json:"profile" gorm:"foreignKey:ProfileId"`
	InternJob   *InternJob `json:"intern_job" gorm:"foreignKey:InternJobId"`

	CreatedAt time.Time      `json:"createdAt" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updatedAt" swaggerignore:"true"`
}
