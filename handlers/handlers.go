package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Client       *mongo.Client
	DatabaseName string
}
type jwtBody struct {
	Token    string
	Username string
}
type paginationResponse struct {
	Data          interface{} `json:"data"`
	Limit         string      `json:"limit"`
	Page          string      `json:"page"`
	TotalElements int64       `json:"totalElements"`
}
type apiResponse struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Error  string      `json:"error"`
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch mux.CurrentRoute(r).GetName() {
		case "CreateRecipe", "ListRecipes", "Renew", "CreateCategory", "GetUserinfo": //Todo better way? Handlers which require auth
			tkn, username, err := lib.VerifyAndReturnJWT(w, r)
			if err != nil {
				log.Println("Error:", err)
				if err := CreateApiResponse(w, nil, http.StatusUnauthorized, "failed", "Failed to authorize user"); err != nil {
					log.Println("Error create Api response:", err)
				}
				return
			}
			ctx := context.WithValue(r.Context(), "token", jwtBody{Token: tkn, Username: username})
			rWithCtx := r.WithContext(ctx)
			h.ServeHTTP(w, rWithCtx)
			return
		}
		h.ServeHTTP(w, r)
	})
}
func CreateApiResponse(w http.ResponseWriter, data interface{}, code int, status string, err string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(apiResponse{
		Data:   data,
		Code:   code,
		Status: status,
		Error:  err,
	}); err != nil {
		return err
	}
	return nil
}
