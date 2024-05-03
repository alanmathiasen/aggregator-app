package controllers

import (
	"fmt"
	"net/http"

	"github.com/alanmathiasen/aggregator-api/auth"
	"github.com/alanmathiasen/aggregator-api/helpers"
	"github.com/alanmathiasen/aggregator-api/services"
	login "github.com/alanmathiasen/aggregator-api/view/auth"
	"github.com/gorilla/sessions"
)

func Login(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		if email == "" {
			session.AddFlash("Please enter your email")
		}
		if password == "" {
			session.AddFlash("Please enter your password")
		}

		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	user, err := services.AuthenticateUser(r.Context(), email, password)
	if err != nil {
		fmt.Print(err.Error())
		session.AddFlash(err.Error())
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}
	session.Values["authenticated"] = true
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func Register(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)

	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		if email == "" {
			session.AddFlash("Please enter your email")
		}
		if password == "" {
			session.AddFlash("Please enter your password")
		}
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/auth/register", http.StatusFound)
		return
	}
	hashedPassword, err := services.HashPassword(password)
	if err != nil {
		session.AddFlash("There was a problem with your registration")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/auth/register", http.StatusFound)
		return
	}
	user, err := services.RegisterUser(r.Context(), email, hashedPassword)
	if err != nil {
		session.AddFlash("This username is already taken")
		fmt.Print(err.Error())
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Print(err.Error())
			return
		}
		http.Redirect(w, r, "/auth/register", http.StatusFound)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func RegisterHTML(w http.ResponseWriter, r *http.Request) {
	component := login.Register()
	err := component.Render(r.Context(), w)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}

func LoginHTML(w http.ResponseWriter, r *http.Request) {
	component := login.Login()
	err := component.Render(r.Context(), w)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}
