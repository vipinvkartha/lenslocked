package controllers

import (
	"fmt"
	"net/http"

	"github.com/vipinvkartha/lenslocked/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User created: %+v", user)
	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, "Could not parse form", http.StatusBadRequest)
	// 	return
	// }
	// // fmt.Fprintf(w, "<p>Email: %s</p>", r.PostForm.Get("email"))
	// // fmt.Fprintf(w, "<p>Password: %s</p>", r.PostForm.Get("password"))
	// fmt.Fprintf(w, "<p>Email: %s</p>", r.FormValue("email"))
	// fmt.Fprintf(w, "<p>Password: %s</p>", r.FormValue("password"))
}
