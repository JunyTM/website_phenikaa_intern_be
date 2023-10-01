package model

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID        uint   `json:"id" gorm:"autoIncrement"`
	UserId    uint   `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Code      string `json:"code"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Birthday  string `json:"birthday"`

	User        *User         `json:"user" gorm:"foreignKey:UserId"`
	Recruitment []Recruitment `json:"recruitment" gorm:"foreignKey:ProfileId"`
	InternShip  []InternShip  `json:"internship" gorm:"foreignKey:ProfileId"`

	CreatedAt time.Time      `json:"createdAt" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updatedAt" swaggerignore:"true"`
}
