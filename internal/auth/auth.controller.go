package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/huynchu/degree-planner-api/config"
	"github.com/huynchu/degree-planner-api/internal/user"
	"github.com/huynchu/degree-planner-api/internal/utils"
)

type AuthController struct {
	userService *user.UserService
}

func NewAuthController(userService *user.UserService) *AuthController {
	InitializeOAuthGoogle()
	return &AuthController{userService: userService}
}

func (c *AuthController) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	HandleLogin(w, r, oauthConfGl, oauthStateStringGl)
}

func (c *AuthController) CallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
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

		// Check if user exists in database
		_, err = c.userService.FindUserByEmail(googleUser.Email)
		if err == user.ErrUserNotFound {
			// If user does not exist, create new user
			_, err = c.userService.CreateNewUser(googleUser.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Load env
		env, err := config.LoadConfig()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create JWT toke
		jwtToken, err := utils.CreateToken(env.TOKEN_EXPIRED_IN, googleUser.Email, env.JWT_SECRET)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    jwtToken,
			Path:     "/",
			MaxAge:   env.TOKEN_MAXAGE * 60,
			Secure:   false, // set to true for https
			HttpOnly: true,
			Domain:   "localhost",
		})

		// Send json response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jwtToken)
	}
}
