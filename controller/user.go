package controller

import (
	"fmt"
	"net/http"

	"github.com/jhampac/gallery/view"
)

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type User struct {
	NewView *view.View
}

func NewUser() *User {
	return &User{
		NewView: view.New("index", "user/new"),
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
	fmt.Fprintln(w, "Email is", form.Email)
	fmt.Fprintln(w, "Password is", form.Password)
}
