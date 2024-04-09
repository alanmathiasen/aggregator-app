package auth

import (
	"encoding/gob"

	"github.com/alanmathiasen/aggregator-api/services"
	"github.com/gorilla/sessions"
)

type keyType string

const SessionKey keyType = "session"

var Store *sessions.CookieStore

func InitStore() {
	Store = sessions.NewCookieStore([]byte("super_duper_secret"))
	Store.Options.HttpOnly = true
	// store.Options.Secure = true
	gob.Register(&services.User{})
}
