package handlers

import (
	"fmt"
	"net/http"

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

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Authentication required")
		fmt.Println("Name:", mux.CurrentRoute(r).GetName())

		h.ServeHTTP(w, r)
	})
}
