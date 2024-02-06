package auth

import (
	"encoding/gob"

	"github.com/alanmathiasen/aggregator-api/services"
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("super_duper_secret"))

func InitStore () {
	Store.Options.HttpOnly = true
	// store.Options.Secure = true
	gob.Register(&services.User{})
}