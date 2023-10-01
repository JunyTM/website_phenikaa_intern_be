package model

import "time"

type UserRole struct {
	ID     uint `json:"id" gorm:"autoIncrement"`
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
	Active bool `json:"active"`

	CreatedAt time.Time `json:"createdAt" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updatedAt" swaggerignore:"true"`
	DeletedAt time.Time `json:"-" swaggerignore:"true"`
}
