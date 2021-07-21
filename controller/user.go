package controller

import (
	"fmt"
	"net/http"

	"github.com/jhampac/gallery/model"
	"github.com/jhampac/gallery/view"
)

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type User struct {
	NewView   *view.View
	LoginView *view.View
	us        *model.UserService
}

func NewUser(us *model.UserService) *User {
	return &User{
		NewView:   view.New("index", "user/new"),
		LoginView: view.New("index", "user/login"),
		us:        us,
	}
}

func (u *User) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := model.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User is", user)
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case model.ErrNotFound:
		fmt.Fprintln(w, "invalid email address")
	case model.ErrInvalidPassword:
		fmt.Fprintln(w, "invalid password")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
