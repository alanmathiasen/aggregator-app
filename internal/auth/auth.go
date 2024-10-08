package auth

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"

	"github.com/alanmathiasen/aggregator-api/internal/services"
	"github.com/gorilla/sessions"
)

type keyType string

const (
	SessionName         = "my-session"
	SessionKey  keyType = "session"
)

var Store *sessions.CookieStore

func InitStore() {
	Store = sessions.NewCookieStore([]byte("super_duper_secret"))
	Store.Options.HttpOnly = true
	// store.Options.Secure = true
	gob.Register(&services.User{})
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, SessionName)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), SessionKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value(SessionKey).(*sessions.Session)
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
