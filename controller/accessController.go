package controller

import (
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

func NewAccessController(db *gorm.DB) AccessController {
	return &accessController{
		accessService: service.NewAccessService(db),
		userService:   service.NewUserService(db),
		db:            db,
	}
}
