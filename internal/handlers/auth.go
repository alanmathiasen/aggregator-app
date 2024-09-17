package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alanmathiasen/aggregator-api/internal/auth"
	"github.com/alanmathiasen/aggregator-api/internal/services"
	login "github.com/alanmathiasen/aggregator-api/internal/views/auth"
	"github.com/alanmathiasen/aggregator-api/pkg/utils"

	"github.com/gorilla/sessions"
)

var authService services.AuthService

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

	user, err := authService.AuthenticateUser(r.Context(), email, password)
	if err != nil {
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
	hashedPassword, err := authService.HashPassword(password)
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
	user, err := authService.RegisterUser(r.Context(), email, hashedPassword)
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

// ****************** PAGES ******************
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)

	var errorMessage string
	for _, f := range session.Flashes() {
		errorMessage += f.(string)
	}

	if err := session.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
	}

	component := login.Register(errorMessage)
	err := component.Render(r.Context(), w)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)

	flashes := session.Flashes()

	var errorMessages []string
	for _, f := range flashes {
		errorMessages = append(errorMessages, f.(string))
	}

	// Save the session to persist the flash clearance
	if err := session.Save(r, w); err != nil {
		log.Printf("LoginPage - Error saving session: %v", err)
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	component := login.Login(strings.Join(errorMessages, " "))
	err := component.Render(r.Context(), w)
	if err != nil {
		log.Printf("LoginPage - Error rendering component: %v", err)
		http.Error(w, "Error rendering component", http.StatusInternalServerError)
		return
	}
}
