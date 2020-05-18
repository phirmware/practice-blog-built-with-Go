package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const (
	pepper = "my-secret-pepper"
)

// UserModel defines the shape of the user
type UserModel struct {
	Username     string
	Email        string
	Password     string
	PasswordHash string
	Remember     string
	RememberHash string
}

// UserService defines the shape of the user
type UserService struct {
	db *gorm.DB
}

// NewUserService returns the UserService
func NewUserService(db *gorm.DB) UserService {
	return UserService{
		db: db,
	}
}

// FindAll finds all users
func (us UserService) FindAll() ([]UserModel, error) {
	var users []UserModel
	if err := us.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Create creates the user in the database
func (us UserService) Create(user *UserModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password+pepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	user.Password = ""
	return us.db.Create(&user).Error
}
