package routing

import (
	"net/http"

	"github.com/Daniorocket/cookit-backend/db"
	"github.com/Daniorocket/cookit-backend/handlers"
	"github.com/Daniorocket/cookit-backend/logger"
	"github.com/gorilla/mux"
)

func NewRouter(
	categoryRepository db.CategoryRepository,
	recipeRepository db.RecipeRepository,
	authRepository db.AuthRepository,
) (*mux.Router, error) {
	h := handlers.Handler{
		CategoryRepository: categoryRepository,
		RecipeRepository:   recipeRepository,
		AuthRepository:     authRepository,
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range initRoutes(h) {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.Use(handlers.AuthenticateMiddleware)
	return router, nil
}
