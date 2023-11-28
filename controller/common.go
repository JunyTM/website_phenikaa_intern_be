package controller

import (
	"context"
	"net/http"
	"phenikaa/infrastructure"
	"phenikaa/model"
	"time"
)

func SaveHttpCookie(fullDomain string, tokenDetail *model.TokenDetail, w http.ResponseWriter) {
	accessCookie := http.Cookie{
		Name:     "AccessToken",
		Domain:   fullDomain,
		Path:     "/",
		Value:    tokenDetail.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * time.Duration(model.AccessTokenTime)),
	}

	refreshCookie := http.Cookie{
		Name:     "RefreshToken",
		Domain:   fullDomain,
		Path:     "/",
		Value:    tokenDetail.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * time.Duration(model.RefreshTokenTime)),
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}

func GetAndDecodeToken(token string) (map[string]interface{}, error) {
	if token == "" {
		return nil, nil
	}
	decodedToken, err := infrastructure.GetDecodeAuth().Decode(token)
	if err != nil {
		return nil, err
	}
	claims, err := decodedToken.AsMap(context.Background())
	if err != nil {
		return nil, err
	}
	return claims, nil
}
