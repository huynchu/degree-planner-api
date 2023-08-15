package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/huynchu/degree-planner-api/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/api/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = ""
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}

func InitializeOAuthGoogle() {
	// load config
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	oauthConfGl.ClientID = env.GOOGLE_CLIENT_ID
	oauthConfGl.ClientSecret = env.GOOGLE_CLIENT_SECRET
	// TODO: generate random string
	oauthStateStringGl = env.OAUTH_STATE_STRING
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	HandleLogin(w, r, oauthConfGl, oauthStateStringGl)
}

func CallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateStringGl {
		http.Error(w, "Invalid Session State", http.StatusUnauthorized)
		return
	}

	code := r.FormValue("code")

	if code == "" {
		http.Error(w, "Invalid Google authorization code", http.StatusUnauthorized)
		return
	} else {
		// Exchange Google auth code for Google tokens
		token, err := oauthConfGl.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Invalid Google authorization code: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Get Google user info from token
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			http.Error(w, "Invalid Google token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		response, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Could not read req body: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Decode response into GoogleUser struct
		var googleUser GoogleUser
		err = json.Unmarshal(response, &googleUser)
		if err != nil {
			http.Error(w, "Could not unmarshal response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(string(response)))
		return
	}
}
