package configs

import (
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"net/http"
)

var SessionStore *sessions.CookieStore

func InitSessionStore() {
	SessionStore = sessions.NewCookieStore([]byte(viper.GetString("sessionSecret")))
	// Make sessions last for 1 week
	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true, // Set to true if serving over HTTPS
	}

}
