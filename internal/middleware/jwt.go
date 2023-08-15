package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/huynchu/degree-planner-api/config"
	"github.com/huynchu/degree-planner-api/internal/user"
	"github.com/huynchu/degree-planner-api/internal/utils"
)

type authCtxKey struct{}

// Auth middleware validates an incoming request jwt token and adds the Account._id of that user to
// the request context. Add this middleware to a route and extract the accountID as follow:
//
//	authedAccountID := ctx.Value(authCtxKey{}).(primitive.ObjectID)

func NewAuthMiddleWare(userService *user.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract access token from the cookie or authorization header
			accessToken := ""
			cookie, err := r.Cookie("access_token")
			if err == nil {
				accessToken = cookie.Value
			} else {
				authHeader := r.Header.Get("Authorization")
				if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
					accessToken = authHeader[7:]
				}
			}

			// Check tokens exists and is not malformed
			if accessToken == "" {
				http.Error(w, "Unauthorized: missing/malformed access token", http.StatusUnauthorized)
				return
			}

			env, err := config.LoadConfig()
			if err != nil {
				http.Error(w, "err loading config", http.StatusInternalServerError)
				return
			}

			claims, err := utils.ValidateToken(accessToken, env.JWT_SECRET)
			if err != nil {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			// Find user in database
			email := claims.(string)
			usr, err := userService.FindUserByEmail(email)
			if err != nil {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			fmt.Println(usr)

			// Add the User to the request context
			ctx := context.WithValue(r.Context(), authCtxKey{}, usr)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
