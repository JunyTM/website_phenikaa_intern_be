package model

import "github.com/lib/pq"

type AdvanceFilterPayload struct {
	ModelType         string         `json:"modelType"`
	IgnoreAssociation bool           `json:"ignoreAssociation"`
	Page              int            `json:"page"`
	PageSize          int            `json:"pageSize"`
	IsPaginateDB      bool           `json:"isPaginateDB"`
	QuerySerch        string         `json:"querySearch"`
	SelectColumn      pq.StringArray `json:"selectColumn"`
}

type BasicQueryPayload struct {
	ModelType string      `json:"modelType"`
	Data      interface{} `json:"data"`
}

type ListModelId struct {
	ID        []uint `gorm:"column:id"`
	ModelType string `json:"modelType"`
}

// TokenDetail details for token authentication
type TokenDetail struct {
	Username     string
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// AccessDetail access detail only from token
type AccessDetail struct {
	AccessUUID string
	UserID     int
}

// Payload for authentication
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"fullName"`
}

type ChangePasswordPayload struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
