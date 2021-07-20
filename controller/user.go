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

type User struct {
	NewView *view.View
	us      *model.UserService
}

func NewUser(us *model.UserService) *User {
	return &User{
		NewView: view.New("index", "user/new"),
		us:      us,
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
		Name:  form.Name,
		Email: form.Email,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User is", user)
}
