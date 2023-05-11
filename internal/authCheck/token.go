package authCheck

import (
	"fmt"
	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

var MySigningKeyUser = []byte("Admin")
var MySigningKeyAdmin = []byte("User")

func Admin(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")
		defer r.Body.Close()

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")

			if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {

				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {

					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return MySigningKeyAdmin, nil
				})

				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					w.Header().Add("Content-Type", "application/json")
					return
				}

				if token.Valid {
					endpoint(w, r)
				}

			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Invalid authorization token")
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Authorization token not found")
		}
	})
}

func UserAndAdmin(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logg.GetLogger()

		w.Header().Set("Connection", "close")
		defer r.Body.Close()

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")

			if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {

				token, err := jwt.Parse(bearerToken[1],
					func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("There was an error")
						}

						return MySigningKeyUser, nil
					})

				token2, err2 := jwt.Parse(bearerToken[1],
					func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("There was an error")
						}

						return MySigningKeyAdmin, nil
					})

				if err != nil && err2 != nil {
					w.WriteHeader(http.StatusForbidden)
					w.Header().Add("Content-Type", "application/json")
					l.Error(err)
					return
				}

				if token.Valid || token2.Valid {
					endpoint(w, r)
				}

			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Invalid authorization token")
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Authorization token not found")
		}
	})
}
