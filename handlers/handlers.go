package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
)

type Handler struct {
	Sess         *mgo.Session
	DatabaseName string
}
type Pagination struct {
	Data          interface{} `json:"data"`
	Limit         int         `json:"limit"`
	Page          int         `json:"page"`
	TotalElements int         `json:"totalElements"`
}
type Token struct {
	Token interface{} `json:"token"`
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch mux.CurrentRoute(r).GetName() {
		case "CreateRecipe", "ListRecipes", "Renew":
			fmt.Println("Auth required!")
			tokenFromRequest := r.Header.Get("Authorization")
			if len(tokenFromRequest) < 10 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			dividedToken := strings.Split(tokenFromRequest, " ")
			tokenFromRequest = dividedToken[1]
			claims := &Claims{}

			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			tkn, err := jwt.ParseWithClaims(tokenFromRequest, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if mux.CurrentRoute(r).GetName() == "Renew" {
				fmt.Println("Jestem w middleware")
				//create a new request context containing the authenticated user
				ctxWithUser := context.WithValue(r.Context(), "token", tokenFromRequest)
				//create a new request using that new context
				rWithUser := r.WithContext(ctxWithUser)
				//call the real handler, passing the new request
				h.ServeHTTP(w, rWithUser)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
