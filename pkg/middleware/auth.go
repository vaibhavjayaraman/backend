package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"historymap-microservices/pkg/tools"
	"net/http"
)

func Auth(unAuthchain http.HandlerFunc) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if validate(w, r) {
				h(w,r)
			} else {
				if unAuthchain != nil {
					unAuthchain(w,r)
				}
			}
		}
	}
}

func validate(w http.ResponseWriter, r *http.Request) bool {
	jwtString := r.Header.Get("Authorization")
	token, err := jwt.Parse(jwtString, func (token *jwt.Token) (interface{}, error){
		return []byte(tools.JwtSecretKey), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

