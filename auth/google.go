package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	"github.com/michaelmagen/task-together/model"

	"github.com/michaelmagen/task-together/configs"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// TODO: Replace redirect uri with real one
var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:3000/auth/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = ""
)

// InitializeOAuthGoogle Function
func InitializeOAuthGoogle() {
	oauthConfGl.ClientID = viper.GetString("google.clientID")
	oauthConfGl.ClientSecret = viper.GetString("google.clientSecret")
	oauthStateStringGl = viper.GetString("oauthStateString")
}

// Login with Google
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	handleLogin(w, r, oauthConfGl, oauthStateStringGl)
}

// TODO: Need to redirect back to front end for errors instead of 400
// Like this -> http.Redirect(w, r, frontendURL+"/login?error=auth_failed", http.StatusFound)
// TODO: Add proper error handeling for this
var redirectString = "/"

func CallbackGoogleOauth(w http.ResponseWriter, r *http.Request) {
	// Ensure that state value is correct
	state := r.FormValue("state")
	if state != oauthStateStringGl {
		log.Info("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Make sure authorization code is present
	code := r.FormValue("code")
	if code == "" {
		log.Warn("Code not found..")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		// User has denied access..
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for token
	token, err := oauthConfGl.Exchange(context.Background(), code)
	if err != nil {
		log.Error("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get user info from google
	user, err := getUserInfo(token)
	if err != nil {
		log.Error("failed to get user info, err:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Info("the user email is:", user.Email, user.UserId)

	// Add user to db
	err = model.CreateUserIfNotExist(user)
	if err != nil {
		log.Error("Failed to add user to db:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Create or update the session with user information
	session, err := configs.SessionStore.Get(r, viper.GetString("sessionName"))
	if err != nil {
		log.Error("Failed to get session:", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	session.Values["userId"] = user.UserId
	session.Values["token"] = token
	// Save the session
	err = session.Save(r, w)
	if err != nil {
		log.Error("Failed to save session: " + err.Error() + "\n")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// TODO: Replace with correct redirect
	http.Redirect(w, r, "/done", http.StatusFound)
}

func getUserInfo(token *oauth2.Token) (*model.User, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = json.Unmarshal(response, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
