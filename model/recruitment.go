package model

import (
	"time"

	"gorm.io/gorm"
)

type Recruitment struct {
	ID          uint `json:"id" gorm:"autoIncrement"`
	ProfileId   uint `json:"profile_id"`
	InternJobId uint `json:"intern_job_id"`

	Profile   *Profile   `json:"profile" gorm:"foreignKey:ProfileId"`
	InternJob *InternJob `json:"intern_job" gorm:"foreignKey:InternJobId"`

	CreatedAt time.Time      `json:"createdAt" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updatedAt" swaggerignore:"true"`
}
