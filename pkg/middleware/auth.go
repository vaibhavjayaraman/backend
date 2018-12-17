package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func AuthMiddleware() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			h(w,r)
		}
	}
}

func validate(w http.ResponseWriter, r *http.Request) {
	jwtString := r.Header.Get("Authorization")
	token, err := jwt.Parse(jwtString, )
}

