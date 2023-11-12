package controller

import (
	"encoding/json"
	"fmt"
	"phenikaa/model"
	"phenikaa/service"

	// "strings"

	"github.com/go-chi/render"
	"gorm.io/gorm"

	"net/http"
)

type accessController struct {
	accessService service.AccessService
	userService   service.UserService
	db            *gorm.DB
}

type AccessController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

// @Summary Login
// @Description Login
// @Tags Access
// @Accept json
// @Produce json
// @Param payload body model.LoginPayload true "Login"
// @Success 200 {object} Response
// @Router /login [post]
func (c *accessController) Login(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var payload model.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if check, err := c.userService.CheckCredentials(payload.Username, payload.Password); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	} else if check != true {
		internalServerErrorResponse(w, r, fmt.Errorf("Credentials was not match, auth: %v", check))
		return
	}

	userInfo, err := c.userService.GetByUsername(payload.Username)
	if err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	tokenDetail, err := c.accessService.CreateToken(userInfo.ID, userInfo.Role)
	if err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	userInfo.AccessToken = tokenDetail.AccessToken
	userInfo.RefreshToken = tokenDetail.RefreshToken

	// fullDomain := r.Header.Get("Origin")
	// fullDomain = strings.Split(fullDomain, "//")[1]
	// SaveHttpCookie(fullDomain, tokenDetail, w)
	res = &Response{
		Data:    userInfo,
		Success: true,
		Message: "Login success",
	}

	render.JSON(w, r, res)
	return
}

// @Summary Logout
// @Description Logout
// @Tags Access
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /logout [post]
func (c *accessController) Logout(w http.ResponseWriter, r *http.Request) {

	return
}

// @Summary Refresh
// @Description Refresh
// @Tags Access
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /refresh [post]
func (c *accessController) Refresh(w http.ResponseWriter, r *http.Request) {
	return
}

func NewAccessController() AccessController {
	return &accessController{
		accessService: service.NewAccessService(),
		userService:   service.NewUserService(),
	}
}
