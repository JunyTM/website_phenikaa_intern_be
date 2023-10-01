package controller

import (
	"encoding/json"
	"net/http"
	"phenikaa/model"

	"github.com/go-chi/render"
)

type UserController interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userService service.UserService
}

// @Summary Create user
// @Description Create user
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body model.User true "User"
// @Success 200 {object} Response
// @Router /user [post]
func (c *userController) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	temp, err := c.userService.Create(user)
	if err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	res := Response{
		Data:    temp,
		Success: true,
		Message: "Create user success",
	}
	render.JSON(w, r, res)
	return
}
