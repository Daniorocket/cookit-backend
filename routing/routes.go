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
			"/api/v1/przepisy",
			false,
			handler.GetListOfRecipes,
		},
		Route{
			"CreateRecipe",
			"POST",
			"/api/v1/dodaj-przepis",
			true,
			handler.CreateRecipe,
		},
		Route{
			"Login",
			"POST",
			"/api/v1/logowanie",
			false,
			handler.Login,
		},
		Route{
			"Register",
			"POST",
			"/api/v1/rejestracja",
			false,
			handler.Register,
		},
		Route{
			"Renew",
			"GET",
			"/api/v1/odnow-token",
			true,
			handler.Renew,
		},
		Route{
			"ListCategories",
			"GET",
			"/api/v1/kategorie",
			false,
			handler.GetListOfCategories,
		},
		Route{
			"CreateCategory",
			"POST",
			"/api/v1/stworz-kategorie",
			true,
			handler.CreateCategory,
		},
		Route{
			"GetCategoryByID",
			"GET",
			"/api/v1/kategorie/{id}",
			false,
			handler.GetCategoryByID,
		},
		Route{
			"GetRecipesByTags",
			"GET",
			"/api/v1/przepisy/kategorie",
			false,
			handler.GetListOfRecipesByCategories,
		},
		Route{
			"GetUserinfo",
			"GET",
			"/api/v1/user",
			true,
			handler.GetUserinfo,
		},
		Route{
			"RemindPassword",
			"POST",
			"/api/v1/przypomnij-haslo",
			false,
			handler.RemindPassword,
		},
		Route{
			"ResetPassword",
			"POST",
			"/api/v1/przypomnij-haslo/{id}",
			false,
			handler.ChangePassword,
		},
		Route{
			"EditUserAccount",
			"PUT",
			"/api/v1/edytuj-konto",
			true,
			handler.EditUserAccount,
		},
	}
}
