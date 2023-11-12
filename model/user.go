package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"type:varchar(100);unique_index"`
	Password string `json:"password"`

	UserRoles []UserRole `json:"user_roles" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"createdAt" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updatedAt" swaggerignore:"true"`
}

type UserResponse struct {
	ID           uint   `json:"id"`
	Role         string `json:"role"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
