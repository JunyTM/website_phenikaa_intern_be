package model

type User struct {
	ID       int    `json:"id" gorm:"primary_key;auto_increment"`
	Username string `json:"username" gorm:"type:varchar(100);unique_index"`
	Password string `json:"password"`
}


