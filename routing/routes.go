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
			"/api/v1/utworz-przepis",
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
			"/api/v1/utworz-kategorie",
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
			"GetUserinfo",
			"GET",
			"/api/v1/informacje-o-uzytkowniku",
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
		Route{
			"GetUnits",
			"GET",
			"/api/v1/jednostki",
			true,
			handler.GetUnits,
		},
		Route{
			"GetRecipeByID",
			"GET",
			"/api/v1/przepis/{id}",
			false,
			handler.GetRecipeByID,
		},
		Route{
			"AddToFavorites",
			"GET",
			"/api/v1/dodaj-do-ulubionych/{id}",
			true,
			handler.AddToFavorites,
		},
		Route{
			"DeleteRecipe",
			"DELETE",
			"/api/v1/usun-przepis/{id}",
			true,
			handler.DeleteRecipe,
		},
		Route{
			"DeleteUserAccount",
			"DELETE",
			"/api/v1/usun-konto",
			true,
			handler.DeleteUserAccount,
		},
		Route{
			"GetUserRecipes",
			"GET",
			"/api/v1/przepisy-uzytkownika",
			true,
			handler.GetUserRecipes,
		},
	}
}
