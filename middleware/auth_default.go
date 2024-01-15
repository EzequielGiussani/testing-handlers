package middleware

import (
	"fmt"
	"net/http"
)

// NewAuthTokenBasic returns a new AuthBasic
func NewAuthTokenBasic(token string) *AuthBasic {
	return &AuthBasic{
		Token: token,
	}
}

// AuthBasic is a struct that contains the basic data of a authenticator
type AuthBasic struct {
	// Token is a string that contains the token
	Token string
}

// Auth is a method that authenticates
func (a *AuthBasic) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Print("Auth before")
		if r.Header.Get("Authorization") != a.Token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
		fmt.Print("Auth after")
	})
}
