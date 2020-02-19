package middleware

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/scalent/ms-boilerplate/cmd/responses"
)

var JWTSecretKey = "secrect"

//ValidateMiddleware to authenticate jwt
func ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(JWTSecretKey), nil
				})
				if error != nil {
					responses.WriteErrorResponse(w, http.StatusBadRequest, "Invalid token")
					return
				}
				if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,x-requested-with, XMLHttpRequest, Access-Control-Allow-Methods")
					next.ServeHTTP(w, r)

				}
			}
		} else {
			responses.WriteErrorResponse(w, http.StatusBadRequest, "An Authorization Header is Required!")
		}
	})
}
