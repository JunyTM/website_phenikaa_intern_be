package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primary_key;auto_increment"`
	Username string `json:"username" gorm:"type:varchar(100);unique_index"`
	Password string `json:"password"`

	UserRoles []UserRole `json:"user_roles" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"createdAt" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updatedAt" swaggerignore:"true"`
}
