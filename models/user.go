package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/hash"
	"lenslocked.com/rand"
)

const (
	pepper        = "my-secret-pepper"
	hmacsecretkey = "my-hmac-key"
)

var (
	// ErrPasswordMissing is returned whnen a password is not provided
	ErrPasswordMissing = errors.New("models: no password was provided")
	// ErrRememberMissing is returned when a remember token is missing
	ErrRememberMissing  = errors.New("models: no token in remember field")
	ErrPasswordTooShort = errors.New("models: Password must ne at least 8 characters")
)

// UserModel defines the shape of the user
type UserModel struct {
	gorm.Model
	Username     string      `gorm:"not null;unique_index"`
	Email        string      `gorm:"not null;unique_index"`
	Password     string      `gorm:"-"`
	PasswordHash string      `gorm:"not null"`
	Remember     string      `gorm:"-"`
	RememberHash string      `gorm:"not null;unique_index"`
	Blogs        []BlogModel `gorm:"-"`
}

// UserDB defines the user db interface
type UserDB interface {
	Create(user *UserModel) error
	FindAll() ([]UserModel, error)
	AutoMigrate() error
}

// UserService defines the shape of the user
type UserService struct {
	UserDB
}

type userGorm struct {
	db *gorm.DB
}

type userValidation struct {
	UserDB
	hmac hash.HMAC
}

type userValFn func(user *UserModel) error

func runUserValFn(user *UserModel, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userGorm{
		db: db,
	}, err
}

func newUserValidation(ug *userGorm) *userValidation {
	hmac := hash.NewHMAC(hmacsecretkey)
	return &userValidation{
		UserDB: ug,
		hmac:   hmac,
	}
}

// NewUserService returns the UserService
func NewUserService(connectionInfo string) (*UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	uv := newUserValidation(ug)
	if err != nil {
		return nil, err
	}
	return &UserService{
		UserDB: uv,
	}, nil
}

// FindAll finds all users
func (ug *userGorm) FindAll() ([]UserModel, error) {
	var users []UserModel
	if err := ug.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Create creates the user in the database
func (ug *userGorm) Create(user *UserModel) error {
	return ug.db.Create(&user).Error
}

func (uv *userValidation) Create(user *UserModel) error {
	if err := runUserValFn(user, uv.checkForPassword, uv.minPasswordLength, uv.hashPassword, uv.generateRemember, uv.hashRemember); err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// ********************* Validation Methods ********************************* //
func (uv *userValidation) hashPassword(user *UserModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password+pepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	user.Password = ""
	return nil
}

func (uv *userValidation) generateRemember(user *UserModel) error {
	randToken, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = randToken
	return nil
}

func (uv *userValidation) hashRemember(user *UserModel) error {
	if user.Remember == "" {
		return ErrRememberMissing
	}
	hash := uv.hmac.Hash(user.Remember)
	user.RememberHash = hash
	return nil
}

func (uv *userValidation) checkForPassword(user *UserModel) error {
	if user.Password == "" {
		return ErrPasswordMissing
	}
	return nil
}

func (uv *userValidation) minPasswordLength(user *UserModel) error {
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

// ********************* Validation Methods ********************************* //

// AutoMigrate creates the user table in the DB
func (ug *userGorm) AutoMigrate() error {
	return ug.db.AutoMigrate(UserModel{}).Error
}
