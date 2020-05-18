package controllers

import (
	"fmt"
	"net/http"

	"practice.blog.com/models"

	"practice.blog.com/views"
)

// User defines the shape of the user
type User struct {
	NewView views.View
	AllView views.View
	us      *models.UserService
}

type signupform struct {
	Username string `schema:"username"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUser returns the user struct
func NewUser(us *models.UserService) User {
	return User{
		NewView: views.NewView("bootstrap", "user/new"),
		AllView: views.NewView("bootstrap", "user/all"),
		us:      us,
	}
}

// New Renders the signup page
func (u User) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Register creates a new user
func (u User) Register(w http.ResponseWriter, r *http.Request) {
	var form signupform
	ParseForm(r, &form)
	user := models.UserModel{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
	}
	u.us.Create(&user)
	fmt.Println(user)
}

// FindAll finds all users
func (u User) FindAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.us.FindAll()
	if err != nil {
		panic(err)
	}
	u.AllView.Render(w, users)
}
