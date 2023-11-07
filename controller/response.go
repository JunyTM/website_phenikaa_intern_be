package controller

import (
	"net/http"
	"phenikaa/model"
	"time"

	"github.com/go-chi/render"
)

type Response struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set(http.StatusText(http.StatusBadRequest), err.Error())
	res := Response{
		Data:    nil,
		Success: false,
		Message: "Bad request. " + err.Error(),
	}
	render.JSON(w, r, res)
}

func internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set(http.StatusText(http.StatusInternalServerError), err.Error())
	res := Response{
		Data:    nil,
		Success: false,
		Message: "Internal Server Error. " + err.Error(),
	}
	render.JSON(w, r, res)
}

func unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set(http.StatusText(http.StatusUnauthorized), err.Error())
	res := Response{
		Data:    nil,
		Success: false,
		Message: "Unauthorized. " + err.Error(),
	}
	render.JSON(w, r, res)
}

func forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set(http.StatusText(http.StatusForbidden), err.Error())
	res := Response{
		Data:    nil,
		Success: false,
		Message: "Forbidden. " + err.Error(),
	}
	render.JSON(w, r, res)
}

func SaveHttpCookie(fullDomain string, tokenDetail *model.TokenDetail, w http.ResponseWriter) {
	accessCookie := http.Cookie{
		Name:   "AccessToken",
		Domain: fullDomain,
		Path:   "/",
		Value:  tokenDetail.AccessToken,
		// HttpOnly: true,
		// Secure:   true,
		Expires: time.Now().Add(time.Hour * time.Duration(model.AccessTokenTime)),
	}

	refreshCookie := http.Cookie{
		Name:   "RefreshToken",
		Domain: fullDomain,
		Path:   "/",
		Value:  tokenDetail.RefreshToken,
		// HttpOnly: true,
		// Secure:   true,
		Expires: time.Now().Add(time.Hour * time.Duration(model.RefreshTokenTime)),
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}
