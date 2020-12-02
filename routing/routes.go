package routing

import (
	"net/http"

	"github.com/Daniorocket/cookit-backend/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	NeedAuth    bool
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func initRoutes(handler handlers.Handler) Routes {
	return Routes{
		Route{
			"ListRecipes",
			"GET",
			"/v1/recipes",
			false,
			handler.ListRecipes,
		},
		Route{
			"CreateRecipe",
			"POST",
			"/v1/recipes",
			true,
			handler.CreateRecipe,
		},
		Route{
			"CreateRecipe",
			"POST",
			"/v1/login",
			false,
			handler.Login,
		},
		Route{
			"CreateRecipe",
			"POST",
			"/v1/register",
			false,
			handler.Register,
		},
		Route{
			"CreateRecipe",
			"POST",
			"/v1/verify",
			false,
			handler.Verify,
		},
	}
}
