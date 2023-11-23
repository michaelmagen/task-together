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

// Gets user information from google api using auth token
func getUserInfo(token *oauth2.Token) (*model.User, error) {
	// Request info from google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read in the response
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarhsall json into a User struct
	var user model.User
	err = json.Unmarshal(response, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Ensures user is authenticated and adds user ID to request context
func AuthedWithGoogle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Find the current session for user
		session, err := configs.SessionStore.Get(r, viper.GetString("sessionName"))
		if err != nil {
			log.Error("Failed to get sessino in auth middleware", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if userId and token exist in the session
		userID, okUserID := session.Values["userId"].(string)
		token, okToken := session.Values["token"].(*oauth2.Token)

		if !okUserID || !okToken {
			log.Error("could not find token or id in session")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Refresh token if need and update session
		newToken, tokenWasRefreshed, err := refreshAuthToken(token)
		if err != nil {
			log.Error("failed to refresh token", err)
			http.Error(w, "Unathorized", http.StatusUnauthorized)
			return
		}

		// If the token is refreshed, then save new token to session
		if tokenWasRefreshed {
			session.Values["token"] = newToken
			err = session.Save(r, w)
			if err != nil {
				log.Error("failed to save session in auth middleware:", err)
				http.Error(w, "Unathorized", http.StatusUnauthorized)
				return
			}
		}

		// If authenticated, store user information in the context for later use
		ctx := context.WithValue(r.Context(), "userId", userID)

		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Attempt to refresh auth token. If token still valid, returns same token. Otherwise refreshes the token
func refreshAuthToken(token *oauth2.Token) (*oauth2.Token, bool, error) {
	isNewToken := false
	// If the token is no longer valid, attempt to refresh the token
	if !token.Valid() {
		// Create token source that can refresh tokens
		tokenSource := oauthConfGl.TokenSource(context.Background(), token)

		// Attempt to refresh the token
		refreshedToken, err := tokenSource.Token()
		if err != nil {
			return nil, false, err
		}

		// Set the token var to the new token
		isNewToken = true
		token = refreshedToken
	}
	return token, isNewToken, nil
}
