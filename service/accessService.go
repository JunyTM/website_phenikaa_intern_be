package service

import (
	"errors"
	"fmt"
	"net/http"
	"phenikaa/infrastructure"
	"phenikaa/model"
	"phenikaa/utils"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/twinj/uuid"
)

type AccessService interface {
	CreateToken(userId uint, role string) (*model.TokenDetail, error)
}

type accessService struct{}

func (s *accessService) CreateToken(userId uint, role string) (*model.TokenDetail, error) {
	var err error
	// Create token details
	tokenDetail := &model.TokenDetail{}

	tokenDetail.AtExpires = time.Now().Add(time.Hour * time.Duration(infrastructure.GetExtendAccessHour())).Unix()
	tokenDetail.AccessUUID = utils.PatternGet(userId) + uuid.NewV4().String()
	tokenDetail.RtExpires = time.Now().Add(time.Hour * time.Duration(infrastructure.GetExtendRefreshHour())).Unix()
	tokenDetail.RefreshUUID = utils.PatternGet(userId) + uuid.NewV4().String()

	// Create Access Token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = tokenDetail.AccessUUID
	atClaims["user_id"] = userId
	atClaims["role"] = role
	atClaims["exp"] = tokenDetail.AtExpires

	_, tokenDetail.AccessToken, err = infrastructure.GetEncodeAuth().Encode(atClaims)
	if err != nil {
		return nil, err
	}

	// Create Resfresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenDetail.RefreshUUID
	rtClaims["user_id"] = userId
	rtClaims["role"] = role
	rtClaims["exp"] = tokenDetail.RtExpires
	_, tokenDetail.RefreshToken, err = infrastructure.GetEncodeAuth().Encode(rtClaims)
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil
}

func (s *accessService) CreateAuth(userID int, tokenDetail *model.TokenDetail) error {
	// converting Unix to UTC(to Time Object)
	accessToken := time.Unix(tokenDetail.AtExpires, 0)
	refreshToken := time.Unix(tokenDetail.RtExpires, 0)
	now := time.Now()

	if errAccess := infrastructure.
		GetRedisClient().
		Set(tokenDetail.AccessUUID, strconv.Itoa(userID), accessToken.Sub(now)).
		Err(); errAccess != nil {
		return errAccess
	}

	if errRefresh := infrastructure.
		GetRedisClient().
		Set(tokenDetail.RefreshUUID, strconv.Itoa(userID), refreshToken.Sub(now)).
		Err(); errRefresh != nil {
		return errRefresh
	}

	return nil
}

func (s *accessService) ExtractTokenMetadata(r *http.Request) (*model.AccessDetail, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return nil, err
	}

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		errLog.Println("can't parse access uuid from token")
		return nil, errors.New("can't parse access uuid from token")
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		errLog.Println(err)
		return nil, err
	}

	return &model.AccessDetail{
		AccessUUID: accessUUID,
		UserID:     int(userID),
	}, nil
}

func NewAccessService() AccessService {
	return &accessService{}
}
