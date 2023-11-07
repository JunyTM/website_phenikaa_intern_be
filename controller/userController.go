package controller

import (
	"encoding/json"
	"net/http"
	"phenikaa/model"
	"phenikaa/service"

	"github.com/go-chi/render"
)

type UserController interface {
	Register(w http.ResponseWriter, r *http.Request)
	ChangePassowrd(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userService service.UserService
}

// @Summary Register
// @Description Register
// @Tags Access
// @Accept json
// @Produce json
// @Param pauload body model.RegisterPayload true "UserRegister"
// @Success 200 {object} Response
// @Router /users/register [post]
func (c *userController) Register(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var payload model.RegisterPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	newUser, err := c.userService.CreateUser(payload)
	if err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	res = &Response{
		Data:    newUser,
		Success: true,
		Message: "Register success",
	}
	render.JSON(w, r, res)
	return
}

// @Summary Change password
// @Description Change password
// @Tags Access
// @Accept json
// @Produce json
// @Param pauload body model.ChangePasswordPayload true "Change password"
// @Success 200 {object} Response
// @Router /users/change-password [put]
func (c *userController) ChangePassowrd(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var payload model.ChangePasswordPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := c.userService.ChangePassword(payload); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	res = &Response{
		Success: true,
		Message: "Change password success",
	}
	render.JSON(w, r, res)
	return
}

// @Summary Reset password
// @Description Reset password
// @Tags Access
// @Accept json
// @Produce json
// @Param pauload body string true "Username"
// @Success 200 {object} Response
// @Router /users/reset-password [put]
func (c *userController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var payload string

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := c.userService.ResetPassword(payload); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	res = &Response{
		Success: true,
		Message: "Reset password success",
	}
	render.JSON(w, r, res)
	return
}

func NewUserController() UserController {
	return &userController{
		userService: service.NewUserService(),
	}
}
