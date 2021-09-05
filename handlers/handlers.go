package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-backend/db"
	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/gorilla/mux"
)

type Handler struct {
	CategoryRepository db.CategoryRepository
	RecipeRepository   db.RecipeRepository
	AuthRepository     db.AuthRepository
}
type jwtBody struct {
	Token    string
	Username string
}
type paginationResponse struct {
	Data          interface{} `json:"data"`
	Limit         int         `json:"limit"`
	Page          int         `json:"page"`
	TotalElements int64       `json:"totalElements"`
}
type apiResponse struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Error  string      `json:"error"`
}

func AuthenticateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch mux.CurrentRoute(r).GetName() {
		case "ListRecipes", "Login", "Register", "ListCategories", "GetCategoryByID", "GetRecipesByCategories", "RemindPassword", "ResetPassword", "GetUnits":
			h.ServeHTTP(w, r)
		default:
			tkn, username, err := lib.VerifyAndReturnJWT(w, r)
			if err != nil {
				log.Println("Error:", err)
				createApiResponse(w, nil, http.StatusUnauthorized, "failed", "Failed to authorize user")
				return
			}
			ctx := context.WithValue(r.Context(), "token", jwtBody{Token: tkn, Username: username})
			rWithCtx := r.WithContext(ctx)
			h.ServeHTTP(w, rWithCtx)
		}
	})
}
func createApiResponse(w http.ResponseWriter, data interface{}, code int, status string, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(apiResponse{
		Data:   data,
		Code:   code,
		Status: status,
		Error:  err,
	}); err != nil {
		log.Println("Error create api response:", err)
	}
}
