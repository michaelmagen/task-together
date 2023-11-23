package configs

import (
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
)

var SessionStore *sessions.CookieStore

func InitSessionStore() {
	SessionStore = sessions.NewCookieStore([]byte(viper.GetString("sessionSecret")))
	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

}
