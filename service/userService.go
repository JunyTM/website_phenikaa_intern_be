package service

import (
	"fmt"
	"log"
	"phenikaa/infrastructure"
	"phenikaa/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserService interface {
	CheckCredentials(username string, password string) (bool, error)
	CheckUserExist(username string) (bool, error)
	Create(newUser model.User) (*model.User, error)
	ResetPassword(username string) error
	ChangePassword(username string, oldPassword string, newPassword string) error
	BanUser(username string) error
}

type userService struct {
	db *gorm.DB
}

func (s *userService) CheckCredentials(username string, password string) (bool, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return comparePassword(user.Password, password), nil
}

func (s *userService) CheckUserExist(username string) (bool, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *userService) Create(newUser model.User) (*model.User, error) {
	var user model.User
	newUser.Password = hashAndSalt(newUser.Password)
	if err := s.db.Model(&user).Clauses(clause.Returning{}).
		Create(&newUser).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&model.UserRole{}).Create(&model.UserRole{
		UserID: user.ID,
		RoleID: 1, // Default role is 1 (client)
	}).Error; err != nil {
		return nil, err
	}
	return &newUser, nil
}

// set default password is 123456
func (s *userService) ResetPassword(username string) error {
	var user model.User
	if err := s.db.Model(&user).Where("username = ?", username).
		Update("password", hashAndSalt("123456")).Error; err != nil {
		return err
	}
	return nil
}

func (s *userService) ChangePassword(username string, oldPassword string, newPassword string) error {
	check, err := s.CheckCredentials(username, oldPassword)
	if err != nil || !check {
		return fmt.Errorf("tai khoan hoac mat khau sai: ", err)
	}

	var user model.User
	if err := s.db.Model(&user).Where("username = ?", username).
		Update("password", hashAndSalt(newPassword)).Error; err != nil {
		return err
	}
	return nil
}

func (s *userService) BanUser(username string) error {
	if err := s.db.Model(&model.User{}).Where("username = ?", username).
		Update("active", false).Error; err != nil {
		return err
	}
	return nil
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePassword(hashedPwd string, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		return false
	}
	return true
}

func NewUserService() UserService {
	return &userService{
		db: infrastructure.GetDB(),
	}
}
