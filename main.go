package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhampac/gallery/controller"
	"github.com/jhampac/gallery/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "password"
	dbname   = "gallery"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	us, err := model.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	userC := controller.NewUser(us)
	staticPage := controller.NewStaticPage()

	r := mux.NewRouter()
	r.Handle("/", staticPage.Home).Methods("GET")
	r.Handle("/contact", staticPage.Contact).Methods("GET")
	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Create).Methods("POST")
	r.Handle("/login", userC.LoginView).Methods("GET")
	r.HandleFunc("/login", userC.Login).Methods("POST")

	r.HandleFunc("/cookietest", userC.CookieTest)

	http.ListenAndServe(":8080", r)
}
