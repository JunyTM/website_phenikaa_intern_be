package service

import (
	"phenikaa/infrastructure"
	"phenikaa/model"
	"phenikaa/utils"

	"time"

	"github.com/myesui/uuid"
	"gorm.io/gorm"
)

type AccessService interface {
	CreateToken(username string) (*model.TokenDetail, error)
}

type accessService struct {
	userService UserService
	db          *gorm.DB
}

func (s *accessService) CreateToken(username string) (*model.TokenDetail, error) {
	tokenDetail := &model.TokenDetail{}
	var err error
	var roleCode string

	var user model.User
	if err = s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	var userRoles model.UserRole
	if err = s.db.Where("user_id = ?", user.ID).Preload("Role").First(&userRoles).Error; err != nil {
		return nil, err
	}
	if userRoles.Role == nil || userRoles.Role.Code == "" {
		roleCode = "client" // Mặc định là người dùng
	} else {
		roleCode = userRoles.Role.Code
	}

	tokenDetail.Username = user.Username
	tokenDetail.AtExpires = time.Now().Add(time.Hour * time.Duration(model.AccessTokenTime)).Unix()
	tokenDetail.AccessUUID = utils.PatternGet(user.ID) + uuid.NewV4().String()
	tokenDetail.RtExpires = time.Now().Add(time.Hour * time.Duration(model.RefreshTokenTime)).Unix()
	tokenDetail.RefreshUUID = utils.PatternGet(user.ID) + uuid.NewV4().String()

	atClaims := make(map[string]interface{})
	atClaims["user_id"] = user.ID
	atClaims["username"] = tokenDetail.Username
	atClaims["access_uuid"] = tokenDetail.AccessUUID
	atClaims["exp"] = tokenDetail.AtExpires
	atClaims["role"] = roleCode
	_, tokenDetail.AccessToken, err = infrastructure.GetEncodeAuth().Encode(atClaims)
	if err != nil {
		return nil, err
	}

	rtClaims := make(map[string]interface{})
	rtClaims["user_id"] = user.ID
	rtClaims["username"] = tokenDetail.Username
	rtClaims["refresh_uuid"] = tokenDetail.RefreshUUID
	rtClaims["exp"] = tokenDetail.RtExpires
	rtClaims["role"] = roleCode
	_, tokenDetail.RefreshToken, err = infrastructure.GetEncodeAuth().Encode(rtClaims)
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil
}

func NewAccessService() AccessService {
	return &accessService{
		userService: NewUserService(),
		db:          infrastructure.GetDB(),
	}
}
