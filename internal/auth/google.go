package auth

import (
	"context"
	"fmt"
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
	fmt.Println(state)
	if state != oauthStateStringGl {
		fmt.Println("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	fmt.Println(code)

	if code == "" {
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		// User has denied access..
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		token, err := oauthConfGl.Exchange(context.Background(), code)
		if err != nil {
			fmt.Println("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			return
		}
		fmt.Println("TOKEN>> AccessToken>> " + token.AccessToken)
		fmt.Println("TOKEN>> Expiration Time>> " + token.Expiry.String())
		fmt.Println("TOKEN>> RefreshToken>> " + token.RefreshToken)

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			fmt.Println("Get: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		response, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ReadAll: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		fmt.Println("parseResponseBody: " + string(response) + "\n")

		w.Write([]byte("Hello, I'm protected\n"))
		w.Write([]byte(string(response)))
		return
	}
}
