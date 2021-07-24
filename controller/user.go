package controller

import (
	"fmt"
	"net/http"

	"github.com/jhampac/gallery/model"
	"github.com/jhampac/gallery/rando"
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

	// get user data from request
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	// create the user
	user := model.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set remember token
	err := u.setRememberCookie(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// redirect
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			fmt.Fprintln(w, "invalid email address")
		case model.ErrInvalidPassword:
			fmt.Fprintln(w, "invalid password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	err = u.setRememberCookie(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *User) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		fmt.Fprintln(w, "Token no longer valid, please login")
		return
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Hello, %s\n", user.Name)
}

func (u *User) setRememberCookie(w http.ResponseWriter, user *model.User) error {
	if user.Remember == "" {
		token, err := rando.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	return nil
}
